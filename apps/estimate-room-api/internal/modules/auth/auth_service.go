package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

type AuthService interface {
	CheckAuth(r *http.Request) (UserID string, err error)
}

var ErrMissingToken = errors.New("missing access token")

type authService struct {
	tokenKey        []byte
	accessTokenRepo repositories.Oauth2AccessTokenRepository
	oidcSessionRepo repositories.Oauth2OidcSessionRepository
	logger          *slog.Logger
}

func NewAuthService(
	tokenKey string,
	accessTokenRepo repositories.Oauth2AccessTokenRepository,
	oidcSessionRepo repositories.Oauth2OidcSessionRepository,
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

	return s.checkToken(token)
}

func (s *authService) checkToken(token string) (UserID string, err error) {
	storedToken, err := s.accessTokenRepo.FindByToken(token)
	if err != nil {
		s.logger.Error("invalid or expired access token")
		return "", errors.New("invalid or expired access token")
	}

	parsedToken, err := utils.ParseToken[repositories.Oauth2AccessTokenModel](s.tokenKey, storedToken.Token)
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
