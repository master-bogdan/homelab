// Package rooms is a module for rooms
package rooms

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
)

type RoomsModule struct {
	Controller RoomsController
	Gateway    *roomsGateway
}

type RoomsModuleDeps struct {
	Router      chi.Router
	WsService   *ws.Service
	AuthService auth.AuthService
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	controller := NewRoomsController()
	gateway := NewRoomsGateway(deps.WsService)

	deps.WsService.Subscribe(EventRoomJoin, gateway.OnEvent)
	deps.WsService.Subscribe(EventRoomLeave, gateway.OnEvent)
	deps.WsService.Subscribe(EventRoomMessage, gateway.OnEvent)

	return &RoomsModule{
		Controller: controller,
		Gateway:    gateway,
	}
}
