package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	emailinfra "github.com/master-bogdan/estimate-room-api/internal/infra/email"
	authdto "github.com/master-bogdan/estimate-room-api/internal/modules/auth/dto"
	authmodels "github.com/master-bogdan/estimate-room-api/internal/modules/auth/models"
	authrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/auth/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	oauth2repositories "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/repositories"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

const passwordResetTokenTTL = time.Hour

var (
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrEmailAlreadyInUse          = errors.New("email already in use")
	ErrInvalidContinueURL         = errors.New("invalid continue url")
	ErrInvalidResetToken          = errors.New("invalid reset token")
	ErrExpiredResetToken          = errors.New("expired reset token")
	ErrUsedResetToken             = errors.New("used reset token")
	ErrGithubAuthNotConfigured    = errors.New("github oauth is not configured")
	ErrGithubAuthenticationFailed = errors.New("github authentication failed")
)

type AuthService interface {
	Login(ctx context.Context, dto *authdto.LoginDTO) (*usersmodels.UserModel, string, error)
	Register(ctx context.Context, dto *authdto.RegisterDTO) (*usersmodels.UserModel, string, error)
	ForgotPassword(ctx context.Context, dto *authdto.ForgotPasswordDTO) error
	ValidateResetPasswordToken(ctx context.Context, token string) (bool, string, error)
	ResetPassword(ctx context.Context, dto *authdto.ResetPasswordDTO) error
	GetSession(r *http.Request) (*usersmodels.UserModel, bool, error)
	Logout(ctx context.Context, sessionID string) error
	StartGithubLogin(continueURL string) (string, error)
	HandleGithubCallback(ctx context.Context, code, stateToken string) (string, string, error)
}

type AuthServiceDeps struct {
	UserService            users.UsersService
	Oauth2Service          oauth2.Oauth2Service
	SessionService         oauth2.Oauth2SessionAuthService
	FrontendBaseURL        string
	EmailClient            emailinfra.Client
	PasswordResetTokenRepo authrepositories.PasswordResetTokenRepository
	AuthCodeRepo           oauth2repositories.Oauth2AuthCodeRepository
	AccessTokenRepo        oauth2repositories.AccessTokenRepository
	RefreshTokenRepo       oauth2repositories.Oauth2RefreshTokenRepository
	OidcSessionRepo        oauth2repositories.OidcSessionRepository
	Github                 oauth2utils.GithubConfig
}

type authService struct {
	userService            users.UsersService
	oauth2Service          oauth2.Oauth2Service
	sessionService         oauth2.Oauth2SessionAuthService
	frontendBaseURL        string
	emailClient            emailinfra.Client
	passwordResetTokenRepo authrepositories.PasswordResetTokenRepository
	authCodeRepo           oauth2repositories.Oauth2AuthCodeRepository
	accessTokenRepo        oauth2repositories.AccessTokenRepository
	refreshTokenRepo       oauth2repositories.Oauth2RefreshTokenRepository
	oidcSessionRepo        oauth2repositories.OidcSessionRepository
	github                 oauth2utils.GithubConfig
	stateTokenKey          []byte
	httpClient             *http.Client
	logger                 *slog.Logger
}

type githubState struct {
	ExpiresAt   int64  `json:"exp"`
	ContinueURL string `json:"continue"`
}

func NewAuthService(deps AuthServiceDeps) AuthService {
	emailClient := deps.EmailClient
	if emailClient == nil {
		emailClient = emailinfra.NewNoopClient()
	}

	return &authService{
		userService:            deps.UserService,
		oauth2Service:          deps.Oauth2Service,
		sessionService:         deps.SessionService,
		frontendBaseURL:        deps.FrontendBaseURL,
		emailClient:            emailClient,
		passwordResetTokenRepo: deps.PasswordResetTokenRepo,
		authCodeRepo:           deps.AuthCodeRepo,
		accessTokenRepo:        deps.AccessTokenRepo,
		refreshTokenRepo:       deps.RefreshTokenRepo,
		oidcSessionRepo:        deps.OidcSessionRepo,
		github:                 deps.Github,
		stateTokenKey:          []byte(deps.Github.StateSecret),
		httpClient:             &http.Client{Timeout: 10 * time.Second},
		logger:                 logger.L().With(slog.String("module", "auth")),
	}
}

func (s *authService) Login(ctx context.Context, dto *authdto.LoginDTO) (*usersmodels.UserModel, string, error) {
	_ = ctx

	email := normalizeEmail(dto.Email)
	query, err := s.parseContinueURL(dto.ContinueURL)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userService.FindByEmail(email)
	if err != nil {
		// TODO: Refactor this
		if errors.Is(err, apperrors.ErrUserNotFound) {
			if hasDeletedEmail, lookupErr := s.userService.HasSoftDeletedEmail(email); lookupErr == nil && hasDeletedEmail {
				return nil, "", ErrInvalidCredentials
			}
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" || !utils.CheckPasswordHash(dto.Password, *user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	sessionID, err := s.createContinueSessionFromQuery(user.UserID, query)
	if err != nil {
		return nil, "", err
	}

	if err := s.userService.UpdateLastLoginAt(user.UserID); err != nil {
		return nil, "", err
	}

	updatedUser, err := s.userService.FindByID(user.UserID)
	if err != nil {
		return nil, "", err
	}

	return updatedUser, sessionID, nil
}

// TODO: Refactor this service
func (s *authService) Register(ctx context.Context, dto *authdto.RegisterDTO) (*usersmodels.UserModel, string, error) {
	_ = ctx

	email := normalizeEmail(dto.Email)

	if _, err := s.userService.FindByEmail(email); err == nil {
		return nil, "", ErrEmailAlreadyInUse
	} else if !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, "", err
	}

	hasDeletedEmail, err := s.userService.HasSoftDeletedEmail(email)
	if err != nil {
		return nil, "", err
	}
	if hasDeletedEmail {
		return nil, "", ErrEmailAlreadyInUse
	}

	query, err := s.parseContinueURL(dto.ContinueURL)
	if err != nil {
		return nil, "", err
	}

	passwordHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, "", err
	}

	displayName := strings.TrimSpace(dto.DisplayName)
	if displayName == "" {
		displayName = defaultDisplayName(email)
	}

	userID, err := s.userService.Create(
		email,
		passwordHash,
		displayName,
		normalizeOptionalProfileField(dto.Organization),
		normalizeOptionalProfileField(dto.Occupation),
	)
	if err != nil {
		return nil, "", err
	}

	sessionID, err := s.createContinueSessionFromQuery(userID, query)
	if err != nil {
		return nil, "", err
	}

	if err := s.userService.UpdateLastLoginAt(userID); err != nil {
		return nil, "", err
	}

	user, err := s.userService.FindByID(userID)
	if err != nil {
		return nil, "", err
	}

	return user, sessionID, nil
}

func (s *authService) ForgotPassword(ctx context.Context, dto *authdto.ForgotPasswordDTO) error {
	normalizedEmail := normalizeEmail(dto.Email)

	user, err := s.userService.FindByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil
		}
		return err
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return nil
	}

	rawToken, err := generateOpaqueToken()
	if err != nil {
		return err
	}

	_, err = s.passwordResetTokenRepo.Create(ctx, &authmodels.PasswordResetTokenModel{
		UserID:    user.UserID,
		TokenHash: hashOpaqueToken(rawToken),
		ExpiresAt: time.Now().Add(passwordResetTokenTTL),
	})
	if err != nil {
		return err
	}

	resetURL, err := s.buildResetPasswordURL(rawToken)
	if err != nil {
		return err
	}

	if err := s.emailClient.Send(ctx, emailinfra.Message{
		To:      []string{normalizedEmail},
		Subject: "Reset your EstimateRoom password",
		TextBody: "We received a request to reset your EstimateRoom password.\n\n" +
			"Use this link to choose a new password:\n" + resetURL + "\n\n" +
			"This link expires in 1 hour. If you did not request a password reset, you can ignore this email.\n",
	}); err != nil {
		return err
	}

	return nil
}

func (s *authService) ValidateResetPasswordToken(ctx context.Context, token string) (bool, string, error) {
	resetToken, err := s.passwordResetTokenRepo.FindByTokenHash(ctx, hashOpaqueToken(token))
	if err != nil {
		if errors.Is(err, apperrors.ErrPasswordResetTokenNotFound) {
			return false, "invalid", nil
		}
		return false, "", err
	}

	switch {
	case resetToken.UsedAt != nil:
		return false, "used", nil
	case time.Now().After(resetToken.ExpiresAt):
		return false, "expired", nil
	default:
		return true, "", nil
	}
}

func (s *authService) ResetPassword(ctx context.Context, dto *authdto.ResetPasswordDTO) error {
	resetToken, err := s.passwordResetTokenRepo.FindByTokenHash(ctx, hashOpaqueToken(dto.Token))
	if err != nil {
		if errors.Is(err, apperrors.ErrPasswordResetTokenNotFound) {
			return ErrInvalidResetToken
		}
		return err
	}

	switch {
	case resetToken.UsedAt != nil:
		return ErrUsedResetToken
	case time.Now().After(resetToken.ExpiresAt):
		return ErrExpiredResetToken
	}

	passwordHash, err := utils.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	if err := s.userService.UpdatePasswordHash(resetToken.UserID, passwordHash); err != nil {
		return err
	}
	if err := s.passwordResetTokenRepo.MarkUsed(ctx, resetToken.PasswordResetTokenID); err != nil {
		return err
	}

	return s.revokeAllUserSessions(ctx, resetToken.UserID)
}

func (s *authService) GetSession(r *http.Request) (*usersmodels.UserModel, bool, error) {
	sessionID := oauth2.ReadOauth2SessionID(r)
	if sessionID == "" {
		return nil, false, nil
	}

	session, err := s.sessionService.GetOidcSessionByID(sessionID)
	if err != nil {
		return nil, false, nil
	}

	user, err := s.userService.FindByID(session.UserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return user, true, nil
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	if strings.TrimSpace(sessionID) == "" {
		return nil
	}

	return s.revokeSession(ctx, sessionID)
}

func (s *authService) StartGithubLogin(continueURL string) (string, error) {
	if !s.github.IsConfigured() {
		return "", ErrGithubAuthNotConfigured
	}

	if _, err := s.parseContinueURL(continueURL); err != nil {
		return "", err
	}

	stateToken, err := utils.GenerateToken(s.stateTokenKey, githubState{
		ExpiresAt:   time.Now().Add(10 * time.Minute).Unix(),
		ContinueURL: continueURL,
	})
	if err != nil {
		return "", err
	}

	return oauth2utils.BuildGithubAuthorizeURL(s.github, stateToken), nil
}

func (s *authService) HandleGithubCallback(ctx context.Context, code, stateToken string) (string, string, error) {
	if !s.github.IsConfigured() {
		return "", "", ErrGithubAuthNotConfigured
	}

	state, err := s.parseGithubState(stateToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := oauth2utils.ExchangeGithubCode(ctx, s.httpClient, s.github, code)
	if err != nil {
		return "", "", ErrGithubAuthenticationFailed
	}

	user, err := oauth2utils.FetchGithubUser(ctx, s.httpClient, accessToken)
	if err != nil {
		return "", "", ErrGithubAuthenticationFailed
	}

	emails, err := oauth2utils.FetchGithubEmails(ctx, s.httpClient, accessToken)
	if err != nil {
		logger.FromContext(ctx, s.logger).Warn(logger.Prefix("MODULE", "AUTH", "GITHUB", "Failed to fetch GitHub emails"), "err", err)
	}

	profile := oauth2utils.BuildGithubProfile(user, emails)
	userID, err := s.oauth2Service.GetOrCreateUserFromGithub(profile)
	if err != nil {
		return "", "", err
	}
	if err := s.userService.UpdateLastLoginAt(userID); err != nil {
		return "", "", err
	}

	sessionID, err := s.createContinueSession(userID, state.ContinueURL)
	if err != nil {
		return "", "", err
	}

	return state.ContinueURL, sessionID, nil
}

func (s *authService) createContinueSession(userID, continueURL string) (string, error) {
	query, err := s.parseContinueURL(continueURL)
	if err != nil {
		return "", err
	}

	return s.createContinueSessionFromQuery(userID, query)
}

func (s *authService) createContinueSessionFromQuery(userID string, query *oauth2dto.Oauth2AuthorizeQueryDTO) (string, error) {
	oidcSessionDTO := &oauth2dto.Oauth2CreateOidcSessionDTO{
		UserID:   userID,
		ClientID: query.ClientID,
		Nonce:    query.Nonce,
	}
	if err := oidcSessionDTO.Validate(); err != nil {
		return "", ErrInvalidContinueURL
	}

	return s.oauth2Service.CreateOidcSession(oidcSessionDTO)
}

func (s *authService) parseContinueURL(continueURL string) (*oauth2dto.Oauth2AuthorizeQueryDTO, error) {
	trimmedContinueURL := strings.TrimSpace(continueURL)

	parsedURL, err := url.Parse(trimmedContinueURL)
	if err != nil {
		s.logContinueURLRejection("parse_failed", nil, nil, err)
		return nil, ErrInvalidContinueURL
	}
	if parsedURL.Scheme != "" || parsedURL.Host != "" {
		s.logContinueURLRejection("absolute_url_not_allowed", parsedURL, nil, errors.New("continue url must be relative"))
		return nil, ErrInvalidContinueURL
	}

	normalizedPath := strings.TrimRight(parsedURL.Path, "/")
	if normalizedPath != "/api/v1/oauth2/authorize" && normalizedPath != "/oauth2/authorize" {
		s.logContinueURLRejection("unsupported_path", parsedURL, nil, errors.New("continue url path is not supported"))
		return nil, ErrInvalidContinueURL
	}

	query := &oauth2dto.Oauth2AuthorizeQueryDTO{
		ClientID:            parsedURL.Query().Get("client_id"),
		RedirectURI:         parsedURL.Query().Get("redirect_uri"),
		ResponseType:        parsedURL.Query().Get("response_type"),
		Scopes:              parsedURL.Query().Get("scopes"),
		State:               parsedURL.Query().Get("state"),
		CodeChallenge:       parsedURL.Query().Get("code_challenge"),
		CodeChallengeMethod: parsedURL.Query().Get("code_challenge_method"),
		Nonce:               parsedURL.Query().Get("nonce"),
	}
	if err := query.Validate(); err != nil {
		s.logContinueURLRejection("invalid_query", parsedURL, query, err)
		return nil, ErrInvalidContinueURL
	}
	if err := s.oauth2Service.ValidateClient(query); err != nil {
		s.logContinueURLRejection("client_validation_failed", parsedURL, query, err)
		return nil, ErrInvalidContinueURL
	}

	return query, nil
}

func (s *authService) logContinueURLRejection(
	reason string,
	parsedURL *url.URL,
	query *oauth2dto.Oauth2AuthorizeQueryDTO,
	err error,
) {
	logArgs := []any{"reason", reason}

	if parsedURL != nil {
		logArgs = append(logArgs,
			"continue_scheme", parsedURL.Scheme,
			"continue_host", parsedURL.Host,
			"continue_path", parsedURL.Path,
		)
	}

	if query != nil {
		logArgs = append(logArgs,
			"client_id", query.ClientID,
			"redirect_uri", query.RedirectURI,
			"response_type", query.ResponseType,
			"scopes", query.Scopes,
			"code_challenge_method", query.CodeChallengeMethod,
		)
	}

	if err != nil {
		logArgs = append(logArgs, "err", err)
	}

	s.logger.Warn("continue url rejected", logArgs...)
}

func (s *authService) parseGithubState(stateToken string) (*githubState, error) {
	if strings.TrimSpace(stateToken) == "" {
		return nil, ErrInvalidContinueURL
	}

	state, err := utils.ParseToken[githubState](s.stateTokenKey, stateToken)
	if err != nil {
		return nil, ErrInvalidContinueURL
	}
	if state.ExpiresAt > 0 && time.Now().Unix() > state.ExpiresAt {
		return nil, ErrInvalidContinueURL
	}
	if _, err := s.parseContinueURL(state.ContinueURL); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *authService) revokeSession(ctx context.Context, sessionID string) error {
	if err := s.authCodeRepo.MarkUsedByOidcSessionID(sessionID); err != nil {
		return err
	}
	if err := s.accessTokenRepo.RevokeByOidcSessionID(ctx, sessionID); err != nil {
		return err
	}
	if err := s.refreshTokenRepo.RevokeByOidcSessionID(ctx, sessionID); err != nil {
		return err
	}
	if err := s.oidcSessionRepo.Revoke(sessionID); err != nil {
		return err
	}

	return nil
}

func (s *authService) revokeAllUserSessions(ctx context.Context, userID string) error {
	if err := s.authCodeRepo.MarkUsedByUserID(userID); err != nil {
		return err
	}
	if err := s.accessTokenRepo.RevokeByUserID(ctx, userID); err != nil {
		return err
	}
	if err := s.refreshTokenRepo.RevokeByUserID(ctx, userID); err != nil {
		return err
	}
	if err := s.oidcSessionRepo.RevokeByUserID(userID); err != nil {
		return err
	}

	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *authService) buildResetPasswordURL(token string) (string, error) {
	trimmedBaseURL := strings.TrimRight(strings.TrimSpace(s.frontendBaseURL), "/")
	if trimmedBaseURL == "" {
		return "", ErrInvalidContinueURL
	}

	resetURL, err := url.Parse(trimmedBaseURL + "/reset-password")
	if err != nil {
		return "", err
	}

	query := resetURL.Query()
	query.Set("token", token)
	resetURL.RawQuery = query.Encode()

	return resetURL.String(), nil
}

func defaultDisplayName(email string) string {
	localPart := strings.TrimSpace(strings.Split(email, "@")[0])
	if localPart == "" {
		return ""
	}

	words := strings.FieldsFunc(localPart, func(r rune) bool {
		return r == '.' || r == '-' || r == '_' || r == '+'
	})
	for idx, word := range words {
		if word == "" {
			continue
		}
		words[idx] = strings.ToUpper(word[:1]) + word[1:]
	}

	return strings.TrimSpace(strings.Join(words, " "))
}

func normalizeOptionalProfileField(value string) *string {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return nil
	}

	return &trimmedValue
}

func generateOpaqueToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashOpaqueToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func formatSessionResponse(user *usersmodels.UserModel) authdto.SessionResponse {
	if user == nil {
		return authdto.SessionResponse{
			Authenticated: false,
			User:          nil,
		}
	}

	return authdto.SessionResponse{
		Authenticated: true,
		User: &authdto.SessionUserResponse{
			ID:           user.UserID,
			Email:        user.Email,
			DisplayName:  user.DisplayName,
			Organization: user.Organization,
			Occupation:   user.Occupation,
			AvatarURL:    user.AvatarURL,
		},
	}
}
