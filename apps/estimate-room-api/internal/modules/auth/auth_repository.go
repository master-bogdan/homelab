package auth

import "github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"

type AccessTokenRepository interface {
	FindByToken(token string) (*models.Oauth2AccessTokenModel, error)
}

type OidcSessionRepository interface {
	FindByID(sessionID string) (*models.OidcSessionModel, error)
}
