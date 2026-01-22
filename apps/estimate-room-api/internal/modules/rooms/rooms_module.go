// Package rooms is a module for rooms
package rooms

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
)

type RoomsModule struct {
	Controller RoomsController
	Gateway    ws.Gateway
}

type RoomsModuleDeps struct {
	Router    chi.Router
	WsManager *ws.Manager
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	controller := NewRoomsController()
	gateway := NewRoomsGateway(deps.WsManager)

	deps.Router.Route("/rooms", func(r chi.Router) {
		r.Get("/{roomID}/ws", gateway.HandleConnection)
	})

	return &RoomsModule{
		Controller: controller,
		Gateway:    gateway,
	}
}
