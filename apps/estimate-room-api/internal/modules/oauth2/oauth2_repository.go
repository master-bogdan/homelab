package oauth2

import (
	"context"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"
)

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
	Create(ctx context.Context, model *models.Oauth2RefreshTokenModel) (string, error)
	FindByToken(ctx context.Context, token string) (*models.Oauth2RefreshTokenModel, error)
	Revoke(ctx context.Context, refreshTokenID string) error
}

type Oauth2AccessTokenRepository interface {
	Create(ctx context.Context, model *models.Oauth2AccessTokenModel) error
	FindByToken(ctx context.Context, token string) (*models.Oauth2AccessTokenModel, error)
	Revoke(ctx context.Context, accessTokenID string) error
}

type UserRepository interface {
	FindByID(userID string) (*models.UserModel, error)
	FindByEmail(email string) (*models.UserModel, error)
	FindByGithubID(githubID string) (*models.UserModel, error)
	Create(email, passwordHash string) (string, error)
	CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error)
	UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error
}
