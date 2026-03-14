// Package ws provides websocket module wiring.
package ws

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

const defaultChannel = "app"
const wsAccessTokenQueryParam = "token"
const wsGuestAccessCookieName = "room_guest_token"

var errQueryAccessTokenNotAllowed = errors.New("websocket access tokens in query params are not allowed")

type WsModule struct {
	Service *Service
}

type WsModuleDeps struct {
	Router         chi.Router
	AuthService    oauth2.AuthService
	TokenKey       string
	Server         PubSub
	OriginPatterns []string
}

func NewWsModule(deps WsModuleDeps) *WsModule {
	service := NewService(deps.Server, defaultChannel)
	service.SetOriginPatterns(deps.OriginPatterns)
	tokenKey := []byte(deps.TokenKey)

	deps.Router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		identity, err := resolveIdentity(r, deps.AuthService, tokenKey)
		if err != nil {
			httputils.WriteResponseError(w, apperrors.CreateHttpError(
				apperrors.ErrUnauthorized,
				apperrors.HttpError{Detail: "unauthorized"},
			))
			return
		}
		service.Connect(w, r, identity)
	})

	return &WsModule{
		Service: service,
	}
}

type guestTokenClaims struct {
	ParticipantID string    `json:"participantId"`
	Role          string    `json:"role"`
	ExpiresAt     time.Time `json:"expiresAt"`
}

func resolveIdentity(r *http.Request, authService oauth2.AuthService, tokenKey []byte) (ConnectIdentity, error) {
	if hasQueryAccessToken(r) {
		return ConnectIdentity{}, errQueryAccessTokenNotAllowed
	}

	userID, authErr := authService.CheckAuth(r)
	if authErr == nil && strings.TrimSpace(userID) != "" {
		return ConnectIdentity{
			Type:   IdentityTypeUser,
			UserID: userID,
		}, nil
	}

	guestToken := readGuestTokenFromCookie(r)
	if guestToken == "" {
		if authErr != nil {
			return ConnectIdentity{}, authErr
		}
		return ConnectIdentity{}, errors.New("missing auth credentials")
	}

	claims, err := utils.ParseToken[guestTokenClaims](tokenKey, guestToken)
	if err != nil {
		return ConnectIdentity{}, err
	}
	if !claims.ExpiresAt.IsZero() && claims.ExpiresAt.Before(time.Now()) {
		return ConnectIdentity{}, errors.New("guest token expired")
	}
	if strings.TrimSpace(claims.Role) != "" && !strings.EqualFold(claims.Role, string(IdentityTypeGuest)) && !strings.EqualFold(claims.Role, "GUEST") {
		return ConnectIdentity{}, errors.New("invalid guest role")
	}
	participantID := strings.TrimSpace(claims.ParticipantID)
	if participantID == "" {
		return ConnectIdentity{}, errors.New("invalid guest participant id")
	}

	return ConnectIdentity{
		Type:          IdentityTypeGuest,
		ParticipantID: participantID,
	}, nil
}

func hasQueryAccessToken(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return false
	}

	return strings.TrimSpace(r.URL.Query().Get(wsAccessTokenQueryParam)) != ""
}

func readGuestTokenFromCookie(r *http.Request) string {
	if r == nil {
		return ""
	}

	cookie, err := r.Cookie(wsGuestAccessCookieName)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(cookie.Value)
}
