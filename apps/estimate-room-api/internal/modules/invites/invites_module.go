// Package invites provides centralized invitation services.
package invites

import (
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	"github.com/uptrace/bun"
)

type InvitesModule struct {
	Service    InvitesService
	Repository invitesrepositories.InvitationRepository
}

type InvitesModuleDeps struct {
	DB       *bun.DB
	TokenKey string
}

func NewInvitesModule(deps InvitesModuleDeps) *InvitesModule {
	repo := invitesrepositories.NewInvitationRepository(deps.DB)
	svc := NewInvitesService(repo, deps.TokenKey)

	return &InvitesModule{
		Service:    svc,
		Repository: repo,
	}
}
