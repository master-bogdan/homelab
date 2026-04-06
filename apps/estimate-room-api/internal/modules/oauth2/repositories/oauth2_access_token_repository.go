package oauth2repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type AccessTokenRepository interface {
	Create(ctx context.Context, model *oauth2models.Oauth2AccessTokenModel) error
	FindByToken(ctx context.Context, token string) (*oauth2models.Oauth2AccessTokenModel, error)
	Revoke(ctx context.Context, accessTokenID string) error
	RevokeByOidcSessionID(ctx context.Context, oidcSessionID string) error
	RevokeByUserID(ctx context.Context, userID string) error
}

type oauth2AccessTokenRepository struct {
	db *bun.DB
}

func NewOauth2AccessTokenRepository(db *bun.DB) *oauth2AccessTokenRepository {
	return &oauth2AccessTokenRepository{db: db}
}

func (r *oauth2AccessTokenRepository) Create(ctx context.Context, model *oauth2models.Oauth2AccessTokenModel) error {
	if model.AccessTokenID == "" {
		model.AccessTokenID = uuid.NewString()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column(
			"access_token_id",
			"user_id",
			"client_id",
			"oidc_session_id",
			"refresh_token_id",
			"scopes",
			"token",
			"issued_at",
			"expires_at",
			"issuer",
			"is_revoked",
		).
		Exec(ctx)
	return err
}

func (r *oauth2AccessTokenRepository) FindByToken(ctx context.Context, token string) (*oauth2models.Oauth2AccessTokenModel, error) {
	model := new(oauth2models.Oauth2AccessTokenModel)
	err := r.db.NewSelect().
		Model(model).
		Where("oat.token = ?", token).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrAccessTokenNotFound
		}
		return nil, err
	}

	return model, nil
}

func (r *oauth2AccessTokenRepository) Revoke(ctx context.Context, accessTokenID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2AccessTokenModel)(nil)).
		Set("is_revoked = ?", true).
		Where("access_token_id = ?", accessTokenID).
		Exec(ctx)
	return err
}

func (r *oauth2AccessTokenRepository) RevokeByOidcSessionID(ctx context.Context, oidcSessionID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2AccessTokenModel)(nil)).
		Set("is_revoked = ?", true).
		Where("oidc_session_id = ?", oidcSessionID).
		Where("is_revoked = ?", false).
		Exec(ctx)
	return err
}

func (r *oauth2AccessTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2AccessTokenModel)(nil)).
		Set("is_revoked = ?", true).
		Where("user_id = ?", userID).
		Where("is_revoked = ?", false).
		Exec(ctx)
	return err
}
