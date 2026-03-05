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
	Controller    RoomsController
	Gateway       *roomsGateway
	Service       RoomsService
	TaskService   RoomsTaskService
	InviteService RoomsInviteService
}

type RoomsModuleDeps struct {
	Router      chi.Router
	DB          *bun.DB
	WsService   *ws.Service
	AuthService oauth2.AuthService
	TokenKey    string
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	roomsRepo := roomsrepositories.NewRoomsRepository(deps.DB)
	taskRepo := roomsrepositories.NewRoomTaskRepository(deps.DB)
	voteRepo := roomsrepositories.NewRoomVoteRepository(deps.DB)
	roundRepo := roomsrepositories.NewRoomTaskRoundRepository(deps.DB)
	participantRepo := roomsrepositories.NewRoomParticipantRepository(deps.DB)
	svc := NewRoomsService(roomsRepo, participantRepo)
	taskSvc := NewRoomsTaskService(roomsRepo, taskRepo, participantRepo)
	inviteSvc := NewRoomsInviteService(roomsRepo, participantRepo, deps.TokenKey)
	ctrl := NewRoomsController(svc, taskSvc, inviteSvc, deps.AuthService)
	gw := NewRoomsGateway(deps.WsService, roomsRepo, participantRepo, taskRepo, voteRepo, roundRepo)

	deps.Router.Route("/rooms", func(r chi.Router) {
		r.Post("/", ctrl.CreateRoom)
		r.Get("/{id}", ctrl.GetRoom)
		r.Post("/{id}/invites/{token}", ctrl.JoinInvite)
		r.Patch("/{id}", ctrl.UpdateRoom)
		r.Route("/{id}/tasks", func(taskRouter chi.Router) {
			taskRouter.Post("/", ctrl.CreateTask)
			taskRouter.Get("/", ctrl.ListTasks)
			taskRouter.Get("/{taskId}", ctrl.GetTask)
			taskRouter.Patch("/{taskId}", ctrl.UpdateTask)
			taskRouter.Delete("/{taskId}", ctrl.DeleteTask)
		})
	})

	deps.WsService.Subscribe(RoomsJoin, gw.handleRoomJoin)
	deps.WsService.Subscribe(RoomsTaskSetCurrent, gw.handleTaskSetCurrent)
	deps.WsService.Subscribe(RoomsVoteCast, gw.handleVoteCast)
	deps.WsService.Subscribe(RoomsVoteReveal, gw.handleVoteReveal)
	deps.WsService.Subscribe(RoomsRoundNext, gw.handleRoundNext)
	deps.WsService.SubscribeDisconnect(gw.handleDisconnect)

	return &RoomsModule{
		Controller:    ctrl,
		Gateway:       gw,
		Service:       svc,
		TaskService:   taskSvc,
		InviteService: inviteSvc,
	}
}
