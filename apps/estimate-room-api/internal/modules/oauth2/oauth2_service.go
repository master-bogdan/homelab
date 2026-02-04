package oauth2

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Oauth2Service interface {
	ValidateClient(dto *oauth2dto.AuthorizeQueryDTO) error
	CreateAuthCode(dto *oauth2dto.CreateOauthCodeDTO) (string, error)
	CreateOidcSession(dto *oauth2dto.CreateOidcSessionDTO) (string, error)
	GetLoggedInUserID(sessionID string) (string, error)
	AuthenticateUser(dto *oauth2dto.UserDTO) (string, error)
	RegisterUser(dto *oauth2dto.UserDTO) (string, error)
	GetOrCreateUserFromGithub(profile oauth2utils.GithubProfile) (string, error)
	GetAuthorizationTokens(ctx context.Context, dto *oauth2dto.GetTokenDTO) (oauth2dto.TokenResponseDTO, error)
	GetRefreshTokens(ctx context.Context, dto *oauth2dto.GetTokenDTO) (oauth2dto.TokenResponseDTO, error)
	GenerateTokenPair(ctx context.Context, userID, clientID, oidcSessionID string, scopes []string) (oauth2dto.TokenResponseDTO, error)
}

type oauth2Service struct {
	clientRepo       Oauth2ClientRepository
	authCodeRepo     Oauth2AuthCodeRepository
	userRepo         UserRepository
	oidcSessionRepo  Oauth2OidcSessionRepository
	refreshTokenRepo Oauth2RefreshTokenRepository
	accessTokenRepo  Oauth2AccessTokenRepository
	tokenKey         []byte
	issuer           string
	logger           *slog.Logger
}

func NewOauth2Service(
	clientRepo Oauth2ClientRepository,
	authCodeRepo Oauth2AuthCodeRepository,
	userRepo UserRepository,
	oidcSessionRepo Oauth2OidcSessionRepository,
	refreshTokenRepo Oauth2RefreshTokenRepository,
	accessTokenRepo Oauth2AccessTokenRepository,
	tokenKey []byte,
	issuer string,
) Oauth2Service {
	log := logger.L().With(slog.String("module", "oauth"))
	return &oauth2Service{
		clientRepo:       clientRepo,
		authCodeRepo:     authCodeRepo,
		userRepo:         userRepo,
		oidcSessionRepo:  oidcSessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		accessTokenRepo:  accessTokenRepo,
		tokenKey:         tokenKey,
		issuer:           issuer,
		logger:           log,
	}
}

func (s *oauth2Service) ValidateClient(dto *oauth2dto.AuthorizeQueryDTO) error {
	s.logger.Info("ValidateClient")
	client, err := s.clientRepo.FindByID(dto.ClientID)
	if err != nil {
		return err
	}

	if !slices.Contains(client.RedirectURIs, dto.RedirectURI) {
		return errors.New("invalid client")
	}

	if !slices.Contains(client.ResponseTypes, dto.ResponseType) {
		return errors.New("invalid response_type")
	}

	if !slices.Contains(client.GrantTypes, "authorization_code") {
		return errors.New("unsupported grant_type")
	}

	for _, scope := range strings.FieldsSeq(dto.Scopes) {
		if !slices.Contains(client.Scopes, scope) {
			return errors.New("invalid scope")
		}
	}

	return nil
}

func (s *oauth2Service) GetLoggedInUserID(sessionID string) (string, error) {
	s.logger.Info("GetLoggedInUserID")
	session, err := s.oidcSessionRepo.FindByID(sessionID)
	if err != nil {
		return "", errors.New("session not found")
	}

	return session.UserID, nil
}

func (s *oauth2Service) CreateAuthCode(dto *oauth2dto.CreateOauthCodeDTO) (string, error) {
	s.logger.Info("CreateAuthCode")
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	code := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)

	authCode := &models.Oauth2AuthCodeModel{
		ClientID:            dto.ClientID,
		UserID:              dto.UserID,
		OidcSessionID:       dto.OidcSessionID,
		Code:                code,
		RedirectURI:         dto.RedirectURI,
		Scopes:              strings.Fields(dto.Scopes),
		CodeChallenge:       dto.CodeChallenge,
		CodeChallengeMethod: dto.CodeChallengeMethod,
		ExpiresAt:           time.Now().Add(5 * time.Minute),
	}

	err = s.authCodeRepo.Create(authCode)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *oauth2Service) CreateOidcSession(dto *oauth2dto.CreateOidcSessionDTO) (string, error) {
	s.logger.Info("CreateOidcSession")
	oidcSession := &models.OidcSessionModel{
		UserID:   dto.UserID,
		ClientID: dto.ClientID,
		Nonce:    dto.Nonce,
	}

	oidcSessionID, err := s.oidcSessionRepo.Create(oidcSession)
	if err != nil {
		return "", err
	}

	return oidcSessionID, nil
}

func (s *oauth2Service) AuthenticateUser(dto *oauth2dto.UserDTO) (string, error) {
	s.logger.Info("AuthenticateUser")
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return "", err
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return "", ErrInvalidCredentials
	}

	if dto.Password == "" {
		return "", ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(dto.Password, *user.PasswordHash) {
		return "", ErrInvalidCredentials
	}

	return user.UserID, nil
}

func (s *oauth2Service) RegisterUser(dto *oauth2dto.UserDTO) (string, error) {
	s.logger.Info("RegisterUser")
	passwordHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return "", err
	}

	userID, err := s.userRepo.Create(dto.Email, passwordHash)
	if err != nil {
		return "", err
	}

	return userID, err
}

func (s *oauth2Service) GetOrCreateUserFromGithub(profile oauth2utils.GithubProfile) (string, error) {
	if profile.ID == "" {
		return "", errors.New("github id is required")
	}

	user, err := s.userRepo.FindByGithubID(profile.ID)
	if err == nil {
		if err := s.userRepo.UpdateGithubProfile(user.UserID, profile.ID, profile.DisplayName, profile.AvatarURL, profile.Email); err != nil {
			return "", err
		}
		return user.UserID, nil
	}

	if !errors.Is(err, repositories.ErrUserNotFound) {
		return "", err
	}

	if profile.Email != nil && *profile.Email != "" {
		userByEmail, err := s.userRepo.FindByEmail(*profile.Email)
		if err == nil {
			if err := s.userRepo.UpdateGithubProfile(userByEmail.UserID, profile.ID, profile.DisplayName, profile.AvatarURL, profile.Email); err != nil {
				return "", err
			}
			return userByEmail.UserID, nil
		}
		if !errors.Is(err, repositories.ErrUserNotFound) {
			return "", err
		}
	}

	return s.userRepo.CreateWithGithub(profile.Email, profile.ID, profile.DisplayName, profile.AvatarURL)
}

func (s *oauth2Service) GetAuthorizationTokens(ctx context.Context, dto *oauth2dto.GetTokenDTO) (oauth2dto.TokenResponseDTO, error) {
	s.logger.Info("GetAuthorizationTokens")
	authCode, err := s.authCodeRepo.FindByCode(dto.Code)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, errors.New("invalid auth code")
	}

	if authCode.ClientID != dto.ClientID {
		return oauth2dto.TokenResponseDTO{}, errors.New("invalid client")
	}

	if authCode.RedirectURI != dto.RedirectURI {
		return oauth2dto.TokenResponseDTO{}, errors.New("invalid redirect_uri")
	}

	if authCode.IsUsed {
		return oauth2dto.TokenResponseDTO{}, errors.New("invalid or used auth code")
	}

	if time.Now().After(authCode.ExpiresAt) {
		return oauth2dto.TokenResponseDTO{}, errors.New("auth code expired")
	}

	if authCode.CodeChallengeMethod != "S256" {
		return oauth2dto.TokenResponseDTO{}, errors.New("unsupported code challenge method")
	}

	hashed := sha256.Sum256([]byte(dto.CodeVerifier))
	encodedVerifier := base64.RawURLEncoding.EncodeToString(hashed[:])

	if encodedVerifier != authCode.CodeChallenge {
		return oauth2dto.TokenResponseDTO{}, errors.New("code_verifier does not match code_challenge")
	}

	authCode.IsUsed = true
	err = s.authCodeRepo.MarkUsed(authCode.AuthCodeID)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	return s.GenerateTokenPair(ctx, authCode.UserID, authCode.ClientID, authCode.OidcSessionID, authCode.Scopes)
}

func (s *oauth2Service) GetRefreshTokens(ctx context.Context, dto *oauth2dto.GetTokenDTO) (oauth2dto.TokenResponseDTO, error) {
	s.logger.Info("GetRefreshTokens")
	refreshToken, err := s.refreshTokenRepo.FindByToken(ctx, dto.RefreshToken)
	if err != nil || refreshToken.IsRevoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return oauth2dto.TokenResponseDTO{}, errors.New("invalid or expired refresh token")
	}

	if refreshToken.ClientID != dto.ClientID {
		return oauth2dto.TokenResponseDTO{}, errors.New("invalid client")
	}

	err = s.refreshTokenRepo.Revoke(ctx, refreshToken.RefreshTokenID)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	return s.GenerateTokenPair(ctx, refreshToken.UserID, refreshToken.ClientID, refreshToken.OidcSessionID, refreshToken.Scopes)
}

func (s *oauth2Service) GenerateTokenPair(ctx context.Context, userID, clientID, oidcSessionID string, scopes []string) (oauth2dto.TokenResponseDTO, error) {
	s.logger.Info("GenerateTokenPair")
	accessTokenDuration := time.Minute * 15     // 15 minutes
	refreshTokenDuration := time.Hour * 24 * 30 // 30 days
	idTokenDuration := time.Minute * 15         // 15 minutes

	refreshTokenPayload := models.Oauth2RefreshTokenModel{
		UserID:        userID,
		ClientID:      clientID,
		OidcSessionID: oidcSessionID,
		Scopes:        scopes,
		IssuedAt:      time.Now(),
		ExpiresAt:     time.Now().Add(refreshTokenDuration),
	}
	refreshToken, err := utils.GenerateToken(s.tokenKey, refreshTokenPayload)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	refreshTokenPayload.Token = refreshToken

	refreshTokenID, err := s.refreshTokenRepo.Create(ctx, &refreshTokenPayload)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	accessTokenPayload := models.Oauth2AccessTokenModel{
		UserID:         userID,
		ClientID:       clientID,
		OidcSessionID:  oidcSessionID,
		RefreshTokenID: &refreshTokenID,
		Scopes:         scopes,
		IssuedAt:       time.Now(),
		ExpiresAt:      time.Now().Add(accessTokenDuration),
		Issuer:         s.issuer,
	}
	accessToken, err := utils.GenerateToken(s.tokenKey, accessTokenPayload)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	accessTokenPayload.Token = accessToken

	err = s.accessTokenRepo.Create(ctx, &accessTokenPayload)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	session, err := s.oidcSessionRepo.FindByID(oidcSessionID)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	now := time.Now()
	idTokenPayload := oauth2dto.IDTokenPayload{
		Issuer:    s.issuer,
		Subject:   userID,
		Audience:  clientID,
		ExpiresAt: now.Add(idTokenDuration).Unix(),
		IssuedAt:  now.Unix(),
		Nonce:     session.Nonce,
	}

	idToken, err := utils.GenerateToken(s.tokenKey, idTokenPayload)
	if err != nil {
		return oauth2dto.TokenResponseDTO{}, err
	}

	return oauth2dto.TokenResponseDTO{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		ExpiresIn:    int(accessTokenDuration.Seconds()),
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}, nil
}
