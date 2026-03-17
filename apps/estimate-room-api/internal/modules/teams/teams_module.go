// Package teams provides teams endpoints.
package teams

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/uptrace/bun"
)

type TeamsModule struct {
	Controller    TeamsController
	Service       TeamsService
	InviteService TeamsInviteService
}

type TeamsModuleDeps struct {
	Router         chi.Router
	DB             *bun.DB
	AuthService    oauth2.AuthService
	UserService    users.UsersService
	InvitesService invites.InvitesService
}

func NewTeamsModule(deps TeamsModuleDeps) *TeamsModule {
	teamRepo := teamsrepositories.NewTeamRepository(deps.DB)
	memberRepo := teamsrepositories.NewTeamMemberRepository(deps.DB)
	invitationRepo := invitesrepositories.NewInvitationRepository(deps.DB)
	svc := NewTeamsService(deps.DB, teamRepo, memberRepo)
	inviteSvc := NewTeamsInviteService(teamRepo, memberRepo, invitationRepo, deps.UserService, deps.InvitesService)
	ctrl := NewTeamsController(svc, inviteSvc, deps.AuthService)

	deps.Router.Route("/teams", func(r chi.Router) {
		r.Post("/", ctrl.CreateTeam)
		r.Get("/", ctrl.ListTeams)
		r.Get("/{id}", ctrl.GetTeam)
		r.Post("/{id}/invites", ctrl.CreateInvites)
		r.Delete("/{id}/members/{userId}", ctrl.RemoveMember)
	})

	return &TeamsModule{
		Controller:    ctrl,
		Service:       svc,
		InviteService: inviteSvc,
	}
}
