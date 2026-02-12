// Package rooms is a module for rooms
package rooms

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
)

type RoomsModule struct {
	Controller RoomsController
	Gateway    *roomsGateway
	Service    RoomsService
}

type RoomsModuleDeps struct {
	Router      chi.Router
	DB          *pgxpool.Pool
	WsService   *ws.Service
	AuthService auth.AuthService
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	roomsRepo := roomsrepositories.NewRoomsRepository(deps.DB)
	svc := NewRoomsService(roomsRepo)
	ctrl := NewRoomsController(svc, deps.AuthService)
	gw := NewRoomsGateway(deps.WsService)

	deps.Router.Route("/rooms", func(r chi.Router) {
		r.Post("/", ctrl.CreateRoom)
	})

	deps.WsService.Subscribe(EventRoomJoin, gw.OnEvent)
	deps.WsService.Subscribe(EventRoomLeave, gw.OnEvent)
	deps.WsService.Subscribe(EventRoomMessage, gw.OnEvent)

	return &RoomsModule{
		Controller: ctrl,
		Gateway:    gw,
	}
}
