package gamification

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/uptrace/bun"
)

type GamificationModule struct {
	Controller GamificationController
	Service    GamificationService
}

type GamificationModuleDeps struct {
	Router      chi.Router
	DB          *bun.DB
	AuthService oauth2.AuthService
	WsService   *ws.Service
}

func NewGamificationModule(deps GamificationModuleDeps) *GamificationModule {
	service := NewGamificationService(deps.DB, newWSRewardNotifier(deps.WsService))
	controller := NewGamificationController(service, deps.AuthService)

	deps.Router.Route("/gamification", func(r chi.Router) {
		r.Get("/me", controller.GetMe)
	})

	return &GamificationModule{
		Controller: controller,
		Service:    service,
	}
}
