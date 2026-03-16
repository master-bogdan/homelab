// Package teams provides teams endpoints.
package teams

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/uptrace/bun"
)

type TeamsModule struct {
	Controller TeamsController
	Service    TeamsService
}

type TeamsModuleDeps struct {
	Router      chi.Router
	DB          *bun.DB
	AuthService oauth2.AuthService
}

func NewTeamsModule(deps TeamsModuleDeps) *TeamsModule {
	teamRepo := teamsrepositories.NewTeamRepository(deps.DB)
	memberRepo := teamsrepositories.NewTeamMemberRepository(deps.DB)
	svc := NewTeamsService(deps.DB, teamRepo, memberRepo)
	ctrl := NewTeamsController(svc, deps.AuthService)

	deps.Router.Route("/teams", func(r chi.Router) {
		r.Post("/", ctrl.CreateTeam)
		r.Get("/", ctrl.ListTeams)
		r.Get("/{id}", ctrl.GetTeam)
		r.Delete("/{id}/members/{userId}", ctrl.RemoveMember)
	})

	return &TeamsModule{
		Controller: ctrl,
		Service:    svc,
	}
}
