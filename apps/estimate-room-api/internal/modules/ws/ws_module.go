// Package ws provides websocket module wiring.
package ws

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

const defaultChannel = "app"

type WsModule struct {
	Service *Service
}

type WsModuleDeps struct {
	Router      chi.Router
	AuthService auth.AuthService
	Server      PubSub
}

func NewWsModule(deps WsModuleDeps) *WsModule {
	service := NewService(deps.Server, defaultChannel)

	deps.Router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		userID, err := deps.AuthService.CheckAuth(r)
		if err != nil {
			utils.WriteResponseError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		service.Connect(w, r, userID)
	})

	return &WsModule{
		Service: service,
	}
}
