// Package rooms is a module for rooms
package rooms

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
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
	VoteService   RoomsVoteService
	ExpiryService RoomsExpiryService
}

type RoomsModuleDeps struct {
	Router         chi.Router
	DB             *bun.DB
	WsService      *ws.Service
	AuthService    oauth2.AuthService
	InvitesService invites.InvitesService
}

func NewRoomsModule(deps RoomsModuleDeps) *RoomsModule {
	roomsRepo := roomsrepositories.NewRoomsRepository(deps.DB)
	taskRepo := roomsrepositories.NewRoomTaskRepository(deps.DB)
	voteRepo := roomsrepositories.NewRoomVoteRepository(deps.DB)
	roundRepo := roomsrepositories.NewRoomTaskRoundRepository(deps.DB)
	participantRepo := roomsrepositories.NewRoomParticipantRepository(deps.DB)
	expirySvc := NewRoomsExpiryService(roomsRepo, deps.WsService)
	svc := NewRoomsService(roomsRepo, participantRepo)
	voteSvc := NewRoomsVoteService(roomsRepo, taskRepo, voteRepo, roundRepo, participantRepo, expirySvc)
	taskSvc := NewRoomsTaskService(roomsRepo, taskRepo, voteSvc, participantRepo, expirySvc)
	ctrl := NewRoomsController(svc, taskSvc, deps.InvitesService, deps.AuthService)
	gw := NewRoomsGateway(deps.WsService, roomsRepo, participantRepo, taskRepo, voteRepo, roundRepo, voteSvc, expirySvc)

	deps.Router.Route("/rooms", func(r chi.Router) {
		r.Post("/", ctrl.CreateRoom)
		r.Get("/{id}", ctrl.GetRoom)
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
	deps.WsService.Subscribe(RoomsTaskFinalize, gw.handleTaskFinalize)
	deps.WsService.SubscribeDisconnect(gw.handleDisconnect)

	return &RoomsModule{
		Controller:    ctrl,
		Gateway:       gw,
		Service:       svc,
		TaskService:   taskSvc,
		VoteService:   voteSvc,
		ExpiryService: expirySvc,
	}
}
