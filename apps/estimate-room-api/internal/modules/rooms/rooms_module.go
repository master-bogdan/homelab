// Package rooms is a module for rooms
package rooms

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/uptrace/bun"
)

type RoomsModule struct {
	Controller RoomsController
	Gateway    *roomsGateway
	Service    RoomsService
}

type RoomsModuleDeps struct {
	Router      chi.Router
	DB          *bun.DB
	WsService   *ws.Service
	AuthService oauth2.AuthService
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	roomsRepo := roomsrepositories.NewRoomsRepository(deps.DB)
	svc := NewRoomsService(roomsRepo)
	ctrl := NewRoomsController(svc, deps.AuthService)
	gw := NewRoomsGateway(deps.WsService)

	deps.Router.Route("/rooms", func(r chi.Router) {
		r.Post("/", ctrl.CreateRoom)
		r.Get("/{id}", ctrl.GetRoom)
		r.Route("/{id}/tasks", func(taskRouter chi.Router) {
			taskRouter.Post("/", ctrl.CreateTask)
			taskRouter.Get("/", ctrl.ListTasks)
			taskRouter.Get("/{taskId}", ctrl.GetTask)
			taskRouter.Patch("/{taskId}", ctrl.UpdateTask)
			taskRouter.Delete("/{taskId}", ctrl.DeleteTask)
		})
	})

	deps.WsService.Subscribe(EventRoomJoin, gw.OnEvent)
	deps.WsService.Subscribe(EventRoomLeave, gw.OnEvent)
	deps.WsService.Subscribe(EventRoomMessage, gw.OnEvent)

	return &RoomsModule{
		Controller: ctrl,
		Gateway:    gw,
	}
}
