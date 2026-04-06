package history

import (
	"github.com/go-chi/chi/v5"
	historyrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/history/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/uptrace/bun"
)

type HistoryModule struct {
	Controller HistoryController
	Service    HistoryService
	Repository historyrepositories.HistoryRepository
}

type HistoryModuleDeps struct {
	Router      chi.Router
	DB          *bun.DB
	AuthService oauth2.Oauth2SessionAuthService
}

func NewHistoryModule(deps HistoryModuleDeps) *HistoryModule {
	repo := historyrepositories.NewHistoryRepository(deps.DB)
	teamRepo := teamsrepositories.NewTeamRepository(deps.DB)
	memberRepo := teamsrepositories.NewTeamMemberRepository(deps.DB)
	svc := NewHistoryService(repo, teamRepo, memberRepo)
	ctrl := NewHistoryController(svc, deps.AuthService)

	deps.Router.Route("/history", func(r chi.Router) {
		r.Get("/me/sessions", ctrl.ListMySessions)
		r.Get("/teams/{id}/sessions", ctrl.ListTeamSessions)
		r.Get("/rooms/{id}/summary", ctrl.GetRoomSummary)
	})

	return &HistoryModule{
		Controller: ctrl,
		Service:    svc,
		Repository: repo,
	}
}
