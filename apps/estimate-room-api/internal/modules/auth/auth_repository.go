package auth

import (
	"context"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"
)

type AccessTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*models.Oauth2AccessTokenModel, error)
}

type OidcSessionRepository interface {
	FindByID(sessionID string) (*models.OidcSessionModel, error)
}
