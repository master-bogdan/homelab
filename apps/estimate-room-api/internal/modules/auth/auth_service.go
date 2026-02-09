package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

type AuthService interface {
	CheckAuth(r *http.Request) (UserID string, err error)
}

var ErrMissingToken = errors.New("missing access token")

type authService struct {
	tokenKey        []byte
	accessTokenRepo AccessTokenRepository
	oidcSessionRepo OidcSessionRepository
	logger          *slog.Logger
}

func NewAuthService(
	tokenKey string,
	accessTokenRepo AccessTokenRepository,
	oidcSessionRepo OidcSessionRepository,
) AuthService {
	log := logger.L().With(slog.String("module", "auth"))
	return &authService{
		tokenKey:        []byte(tokenKey),
		accessTokenRepo: accessTokenRepo,
		oidcSessionRepo: oidcSessionRepo,
		logger:          log,
	}
}

func (s *authService) CheckAuth(r *http.Request) (UserID string, err error) {
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

func (s *authService) checkToken(ctx context.Context, token string) (UserID string, err error) {
	storedToken, err := s.accessTokenRepo.FindByToken(ctx, token)
	if err != nil {
		s.logger.Error("invalid or expired access token")
		return "", errors.New("invalid or expired access token")
	}

	parsedToken, err := utils.ParseToken[models.Oauth2AccessTokenModel](s.tokenKey, storedToken.Token)
	if err != nil {
		return "", errors.New("invalid or expired access token")
	}

	if storedToken.IsRevoked || storedToken.ExpiresAt.Before(time.Now()) {
		return "", errors.New("access token is revoked or expired")
	}

	if parsedToken.OidcSessionID == "" {
		return "", errors.New("OIDC session not found in access token")
	}

	session, err := s.oidcSessionRepo.FindByID(parsedToken.OidcSessionID)
	if err != nil {
		return "", errors.New("OIDC session not found")
	}

	if session.UserID != parsedToken.UserID || session.ClientID != parsedToken.ClientID {
		return "", errors.New("OIDC session mismatch")
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
