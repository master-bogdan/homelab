// package oauth2
//
// import (
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/base64"
// 	"errors"
// 	"log/slog"
// 	"slices"
// 	"strings"
// 	"time"
//
// 	"github.com/master-bogdan/clear-cash-api/internal/infra/db/postgresql/repositories"
// 	oauth2_dto "github.com/master-bogdan/clear-cash-api/internal/modules/oauth2/dto"
// 	"github.com/master-bogdan/clear-cash-api/internal/pkg/utils"
// )
//
// type Oauth2Service interface {
// 	ValidateClient(dto *oauth2_dto.AuthorizeQueryDTO) error
// 	CreateAuthCode(dto *oauth2_dto.CreateOauthCodeDTO) (string, error)
// 	CreateOidcSession(dto *oauth2_dto.CreateOidcSessionDTO) (string, error)
// 	GetLoggedInUserID(sessionID string) (string, error)
// 	AuthenticateUser(dto *oauth2_dto.UserDTO) (string, error)
// 	RegisterUser(dto *oauth2_dto.UserDTO) (string, error)
// 	GetAuthorizationTokens(dto *oauth2_dto.GetTokenDTO) (oauth2_dto.TokenResponseDTO, error)
// 	GetRefreshTokens(dto *oauth2_dto.GetTokenDTO) (oauth2_dto.TokenResponseDTO, error)
// 	GenerateTokenPair(userID, clientID, oidcSessionID string, scopes []string) (oauth2_dto.TokenResponseDTO, error)
// }
//
// type oauth2Service struct {
// 	clientRepo       repositories.Oauth2ClientRepository
// 	authCodeRepo     repositories.Oauth2AuthCodeRepository
// 	userRepo         repositories.Oauth2UserRepository
// 	oidcSessionRepo  repositories.Oauth2OidcSessionRepository
// 	refreshTokenRepo repositories.Oauth2RefreshTokenRepository
// 	accessTokenRepo  repositories.Oauth2AccessTokenRepository
// 	tokenKey         []byte
// 	logger           *slog.Logger
// }
//
// func NewOauth2Service(
// 	clientRepo repositories.Oauth2ClientRepository,
// 	authCodeRepo repositories.Oauth2AuthCodeRepository,
// 	userRepo repositories.Oauth2UserRepository,
// 	oidcSessionRepo repositories.Oauth2OidcSessionRepository,
// 	refreshTokenRepo repositories.Oauth2RefreshTokenRepository,
// 	accessTokenRepo repositories.Oauth2AccessTokenRepository,
// 	tokenKey []byte,
// 	logger *slog.Logger,
// ) Oauth2Service {
// 	return &oauth2Service{
// 		clientRepo:       clientRepo,
// 		authCodeRepo:     authCodeRepo,
// 		userRepo:         userRepo,
// 		oidcSessionRepo:  oidcSessionRepo,
// 		refreshTokenRepo: refreshTokenRepo,
// 		accessTokenRepo:  accessTokenRepo,
// 		tokenKey:         tokenKey,
// 		logger:           logger,
// 	}
// }
//
// func (s *oauth2Service) ValidateClient(dto *oauth2_dto.AuthorizeQueryDTO) error {
// 	s.logger.Info("ValidateClient")
// 	client, err := s.clientRepo.FindByID(dto.ClientID)
// 	if err != nil {
// 		return err
// 	}
//
// 	result := slices.Contains(client.RedirectURIs, dto.RedirectURI)
//
// 	if !result {
// 		return errors.New("Invalid client")
// 	}
//
// 	return nil
// }
//
// func (s *oauth2Service) GetLoggedInUserID(sessionID string) (string, error) {
// 	s.logger.Info("GetLoggedInUserID")
// 	session, err := s.oidcSessionRepo.FindByID(sessionID)
// 	if err != nil {
// 		return "", errors.New("Session not found")
// 	}
//
// 	return session.UserID, nil
// }
//
// func (s *oauth2Service) CreateAuthCode(dto *oauth2_dto.CreateOauthCodeDTO) (string, error) {
// 	s.logger.Info("CreateAuthCode")
// 	b := make([]byte, 32)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	code := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
//
// 	authCode := &repositories.Oauth2AuthCodeModel{
// 		ClientID:            dto.ClientID,
// 		UserID:              dto.UserID,
// 		OidcSessionID:       dto.OidcSessionID,
// 		Code:                code,
// 		RedirectURI:         dto.RedirectURI,
// 		Scopes:              strings.Fields(dto.Scopes),
// 		CodeChallenge:       dto.CodeChallenge,
// 		CodeChallengeMethod: dto.CodeChallengeMethod,
// 	}
//
// 	err = s.authCodeRepo.Create(authCode)
// 	if err != nil {
// 		return "", err
// 	}
// 	return code, nil
// }
//
// func (s *oauth2Service) CreateOidcSession(dto *oauth2_dto.CreateOidcSessionDTO) (string, error) {
// 	s.logger.Info("CreateOidcSession")
// 	oidcSession := &repositories.OidcSessionModel{
// 		UserID:   dto.UserID,
// 		ClientID: dto.ClientID,
// 		Nonce:    dto.Nonce,
// 	}
//
// 	oidcSessionID, err := s.oidcSessionRepo.Create(oidcSession)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return oidcSessionID, nil
// }
//
// func (s *oauth2Service) AuthenticateUser(dto *oauth2_dto.UserDTO) (string, error) {
// 	s.logger.Info("AuthenticateUser")
// 	user, err := s.userRepo.FindByEmail(dto.Email)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return user.UserID, nil
// }
//
// func (s *oauth2Service) RegisterUser(dto *oauth2_dto.UserDTO) (string, error) {
// 	s.logger.Info("RegisterUser")
// 	passwordHash, err := utils.HashPassword(dto.Password)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	userID, err := s.userRepo.Create(dto.Email, passwordHash)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return userID, err
// }
//
// func (s *oauth2Service) GetAuthorizationTokens(dto *oauth2_dto.GetTokenDTO) (oauth2_dto.TokenResponseDTO, error) {
// 	s.logger.Info("GetAuthorizationTokens")
// 	authCode, err := s.authCodeRepo.FindByCode(dto.Code)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, errors.New("invalid auth code")
// 	}
//
// 	if authCode.IsUsed {
// 		return oauth2_dto.TokenResponseDTO{}, errors.New("invalid or used auth code")
// 	}
//
// 	if authCode.CodeChallengeMethod != "S256" {
// 		return oauth2_dto.TokenResponseDTO{}, errors.New("unsupported code challenge method")
// 	}
//
// 	hashed := sha256.Sum256([]byte(dto.CodeVerifier))
// 	encodedVerifier := base64.RawURLEncoding.EncodeToString(hashed[:])
//
// 	if encodedVerifier != authCode.CodeChallenge {
// 		return oauth2_dto.TokenResponseDTO{}, errors.New("code_verifier does not match code_challenge")
// 	}
//
// 	authCode.IsUsed = true
// 	err = s.authCodeRepo.MarkUsed(authCode.AuthCodeID)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	return s.GenerateTokenPair(authCode.UserID, authCode.ClientID, authCode.OidcSessionID, authCode.Scopes)
// }
//
// func (s *oauth2Service) GetRefreshTokens(dto *oauth2_dto.GetTokenDTO) (oauth2_dto.TokenResponseDTO, error) {
// 	s.logger.Info("GetRefreshTokens")
// 	refreshToken, err := s.refreshTokenRepo.FindByToken(dto.RefreshToken)
// 	if err != nil || refreshToken.IsRevoked || refreshToken.ExpiresAt.Before(time.Now()) {
// 		return oauth2_dto.TokenResponseDTO{}, errors.New("invalid or expired refresh token")
// 	}
//
// 	err = s.refreshTokenRepo.Revoke(refreshToken.RefreshTokenID)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	return s.GenerateTokenPair(refreshToken.UserID, refreshToken.ClientID, refreshToken.OidcSessionID, refreshToken.Scopes)
// }
//
// func (s *oauth2Service) GenerateTokenPair(userID, clientID, oidcSessionID string, scopes []string) (oauth2_dto.TokenResponseDTO, error) {
// 	s.logger.Info("GenerateTokenPair")
// 	accessTokenDuration := time.Minute * 15     // 15 minutes
// 	refreshTokenDuration := time.Hour * 24 * 30 // 30 days
// 	idTokenDuration := time.Minute * 15         // 15 minutes
//
// 	refreshTokenPayload := repositories.Oauth2RefreshTokenModel{
// 		UserID:        userID,
// 		ClientID:      clientID,
// 		OidcSessionID: oidcSessionID,
// 		Scopes:        scopes,
// 		IssuedAt:      time.Now(),
// 		ExpiresAt:     time.Now().Add(refreshTokenDuration),
// 	}
// 	refreshToken, err := utils.GenerateToken(s.tokenKey, refreshTokenPayload)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	refreshTokenPayload.Token = refreshToken
//
// 	refreshTokenID, err := s.refreshTokenRepo.Create(&refreshTokenPayload)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	accessTokenPayload := repositories.Oauth2AccessTokenModel{
// 		UserID:         userID,
// 		ClientID:       clientID,
// 		OidcSessionID:  oidcSessionID,
// 		RefreshTokenID: &refreshTokenID,
// 		Scopes:         scopes,
// 		IssuedAt:       time.Now(),
// 		ExpiresAt:      time.Now().Add(accessTokenDuration),
// 		// TODO: update this to be dynamic
// 		Issuer: "http://localhost:8000",
// 	}
// 	accessToken, err := utils.GenerateToken(s.tokenKey, accessTokenPayload)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	accessTokenPayload.Token = accessToken
//
// 	err = s.accessTokenRepo.Create(&accessTokenPayload)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	session, err := s.oidcSessionRepo.FindByID(oidcSessionID)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	now := time.Now()
// 	idTokenPayload := oauth2_dto.IDTokenPayload{
// 		// TODO: update this to be dynamic
// 		Issuer:    "http://localhost:8000",
// 		Subject:   userID,
// 		Audience:  clientID,
// 		ExpiresAt: now.Add(idTokenDuration).Unix(),
// 		IssuedAt:  now.Unix(),
// 		Nonce:     session.Nonce,
// 	}
//
// 	idToken, err := utils.GenerateToken(s.tokenKey, idTokenPayload)
// 	if err != nil {
// 		return oauth2_dto.TokenResponseDTO{}, err
// 	}
//
// 	return oauth2_dto.TokenResponseDTO{
// 		AccessToken:  accessToken,
// 		TokenType:    "bearer",
// 		ExpiresIn:    int(accessTokenDuration.Seconds()),
// 		RefreshToken: refreshToken,
// 		IDToken:      idToken,
// 	}, nil
// }
