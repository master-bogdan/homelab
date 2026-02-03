package oauth2

import "github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"

type Oauth2ClientRepository interface {
	FindByID(clientID string) (*models.Oauth2ClientModel, error)
}

type Oauth2AuthCodeRepository interface {
	Create(model *models.Oauth2AuthCodeModel) error
	FindByCode(code string) (*models.Oauth2AuthCodeModel, error)
	MarkUsed(authCodeID string) error
}

type Oauth2OidcSessionRepository interface {
	Create(model *models.OidcSessionModel) (string, error)
	FindByID(sessionID string) (*models.OidcSessionModel, error)
}

type Oauth2RefreshTokenRepository interface {
	Create(model *models.Oauth2RefreshTokenModel) (string, error)
	FindByToken(token string) (*models.Oauth2RefreshTokenModel, error)
	Revoke(refreshTokenID string) error
}

type Oauth2AccessTokenRepository interface {
	Create(model *models.Oauth2AccessTokenModel) error
	FindByToken(token string) (*models.Oauth2AccessTokenModel, error)
	Revoke(accessTokenID string) error
}

type UserRepository interface {
	FindByID(userID string) (*models.UserModel, error)
	FindByEmail(email string) (*models.UserModel, error)
	Create(email, passwordHash string) (string, error)
}
