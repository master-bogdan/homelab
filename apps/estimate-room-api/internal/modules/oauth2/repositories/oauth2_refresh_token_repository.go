package oauth2repositories

import (
	"context"
	"errors"

	"database/sql"
	"github.com/google/uuid"
	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type Oauth2RefreshTokenRepository interface {
	Create(ctx context.Context, model *oauth2models.Oauth2RefreshTokenModel) (string, error)
	FindByToken(ctx context.Context, token string) (*oauth2models.Oauth2RefreshTokenModel, error)
	Revoke(ctx context.Context, refreshTokenID string) error
	RevokeByOidcSessionID(ctx context.Context, oidcSessionID string) error
	RevokeByUserID(ctx context.Context, userID string) error
}

type oauth2RefreshTokenRepository struct {
	db *bun.DB
}

func NewOauth2RefreshTokenRepository(db *bun.DB) *oauth2RefreshTokenRepository {
	return &oauth2RefreshTokenRepository{db: db}
}

func (r *oauth2RefreshTokenRepository) Create(ctx context.Context, model *oauth2models.Oauth2RefreshTokenModel) (string, error) {
	if model.RefreshTokenID == "" {
		model.RefreshTokenID = uuid.NewString()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column(
			"refresh_token_id",
			"user_id",
			"client_id",
			"oidc_session_id",
			"scopes",
			"token",
			"issued_at",
			"expires_at",
			"is_revoked",
		).
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return model.RefreshTokenID, nil
}

func (r *oauth2RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*oauth2models.Oauth2RefreshTokenModel, error) {
	model := new(oauth2models.Oauth2RefreshTokenModel)
	err := r.db.NewSelect().
		Model(model).
		Where("ort.token = ?", token).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return model, nil
}

func (r *oauth2RefreshTokenRepository) Revoke(ctx context.Context, refreshTokenID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2RefreshTokenModel)(nil)).
		Set("is_revoked = ?", true).
		Where("refresh_token_id = ?", refreshTokenID).
		Exec(ctx)
	return err
}

func (r *oauth2RefreshTokenRepository) RevokeByOidcSessionID(ctx context.Context, oidcSessionID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2RefreshTokenModel)(nil)).
		Set("is_revoked = ?", true).
		Where("oidc_session_id = ?", oidcSessionID).
		Where("is_revoked = ?", false).
		Exec(ctx)
	return err
}

func (r *oauth2RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2RefreshTokenModel)(nil)).
		Set("is_revoked = ?", true).
		Where("user_id = ?", userID).
		Where("is_revoked = ?", false).
		Exec(ctx)
	return err
}
