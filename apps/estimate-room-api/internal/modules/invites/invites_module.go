// Package invites provides centralized invitation services.
package invites

import (
	"github.com/go-chi/chi/v5"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	"github.com/uptrace/bun"
)

type InvitesModule struct {
	Controller InvitesController
	Service    InvitesService
	Repository invitesrepositories.InvitationRepository
}

type InvitesModuleDeps struct {
	Router      chi.Router
	DB          *bun.DB
	AuthService oauth2.Oauth2SessionAuthService
	TokenKey    string
}

func NewInvitesModule(deps InvitesModuleDeps) *InvitesModule {
	repo := invitesrepositories.NewInvitationRepository(deps.DB)
	svc := NewInvitesService(deps.DB, repo, deps.TokenKey)
	ctrl := NewInvitesController(svc, deps.AuthService)

	deps.Router.Route("/invites", func(r chi.Router) {
		r.Get("/{token}", ctrl.PreviewInvitation)
		r.Post("/{token}/accept", ctrl.AcceptInvitation)
		r.Post("/{token}/decline", ctrl.DeclineInvitation)
		r.Post("/{id}/revoke", ctrl.RevokeInvitation)
	})

	return &InvitesModule{
		Controller: ctrl,
		Service:    svc,
		Repository: repo,
	}
}
