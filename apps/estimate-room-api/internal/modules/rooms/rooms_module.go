// Package rooms is a module for rooms
package rooms

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
)

type RoomsModule struct {
	Controller RoomsController
	Gateway    *roomsGateway
}

type RoomsModuleDeps struct {
	Router      chi.Router
	WsManager   *ws.Manager
	AuthService auth.AuthService
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	controller := NewRoomsController()
	gateway := NewRoomsGateway(deps.WsManager)

	deps.WsManager.Subscribe(EventRoomJoin, gateway.OnEvent)
	deps.WsManager.Subscribe(EventRoomLeave, gateway.OnEvent)
	deps.WsManager.Subscribe(EventRoomMessage, gateway.OnEvent)

	return &RoomsModule{
		Controller: controller,
		Gateway:    gateway,
	}
}
