package oauth2

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
	oauth2repositories "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
	"github.com/uptrace/bun"
)

type AuthService interface {
	CheckAuth(r *http.Request) (userID string, err error)
	CreateOidcSession(model *oauth2models.OidcSessionModel) (string, error)
	GetOidcSessionByID(sessionID string) (*oauth2models.OidcSessionModel, error)
	CreateAccessToken(ctx context.Context, model *oauth2models.Oauth2AccessTokenModel) error
}

var ErrMissingToken = errors.New("missing access token")

type authService struct {
	tokenKey        []byte
	accessTokenRepo oauth2repositories.AccessTokenRepository
	oidcSessionRepo oauth2repositories.OidcSessionRepository
	logger          *slog.Logger
}

func NewAuthService(
	tokenKey string,
	accessTokenRepo oauth2repositories.AccessTokenRepository,
	oidcSessionRepo oauth2repositories.OidcSessionRepository,
) AuthService {
	log := logger.L().With(slog.String("module", "oauth2-auth"))
	return &authService{
		tokenKey:        []byte(tokenKey),
		accessTokenRepo: accessTokenRepo,
		oidcSessionRepo: oidcSessionRepo,
		logger:          log,
	}
}

func NewAuthServiceFromDB(tokenKey string, db *bun.DB) AuthService {
	accessTokenRepo := oauth2repositories.NewOauth2AccessTokenRepository(db)
	oidcSessionRepo := oauth2repositories.NewOauth2OidcSessionRepository(db)
	return NewAuthService(tokenKey, accessTokenRepo, oidcSessionRepo)
}

func (s *authService) CheckAuth(r *http.Request) (string, error) {
	token := extractToken(r)
	if token == "" {
		return "", ErrMissingToken
	}

	ctx := context.Background()
	if r != nil {
		ctx = r.Context()
	}

	return s.checkToken(ctx, token)
}

func (s *authService) CreateOidcSession(model *oauth2models.OidcSessionModel) (string, error) {
	return s.oidcSessionRepo.Create(model)
}

func (s *authService) GetOidcSessionByID(sessionID string) (*oauth2models.OidcSessionModel, error) {
	return s.oidcSessionRepo.FindByID(sessionID)
}

func (s *authService) CreateAccessToken(ctx context.Context, model *oauth2models.Oauth2AccessTokenModel) error {
	return s.accessTokenRepo.Create(ctx, model)
}

func (s *authService) checkToken(ctx context.Context, token string) (string, error) {
	storedToken, err := s.accessTokenRepo.FindByToken(ctx, token)
	if err != nil {
		s.logger.Error("invalid or expired access token")
		return "", errors.New("invalid or expired access token")
	}

	parsedToken, err := utils.ParseToken[oauth2models.Oauth2AccessTokenModel](s.tokenKey, storedToken.Token)
	if err != nil {
		return "", errors.New("invalid or expired access token")
	}

	if storedToken.IsRevoked || storedToken.ExpiresAt.Before(time.Now()) {
		return "", errors.New("access token is revoked or expired")
	}

	if parsedToken.OidcSessionID == "" {
		return "", errors.New("oidc session not found in access token")
	}

	session, err := s.oidcSessionRepo.FindByID(parsedToken.OidcSessionID)
	if err != nil {
		return "", errors.New("oidc session not found")
	}

	if session.UserID != parsedToken.UserID || session.ClientID != parsedToken.ClientID {
		return "", errors.New("oidc session mismatch")
	}

	return parsedToken.UserID, nil
}

func extractToken(r *http.Request) string {
	if r == nil {
		return ""
	}

	if token := bearerToken(r.Header.Get("Authorization")); token != "" {
		return token
	}

	cookie, err := r.Cookie("access_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

func bearerToken(value string) string {
	if value == "" {
		return ""
	}

	scheme, token, ok := strings.Cut(value, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}

	return strings.TrimSpace(token)
}
