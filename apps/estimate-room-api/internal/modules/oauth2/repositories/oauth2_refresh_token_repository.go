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
	const query = `
		INSERT INTO oauth2_refresh_tokens (
			refresh_token_id, user_id, client_id, oidc_session_id, scopes,
			token, issued_at, expires_at, is_revoked
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if model.RefreshTokenID == "" {
		model.RefreshTokenID = uuid.NewString()
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		model.RefreshTokenID,
		model.UserID,
		model.ClientID,
		model.OidcSessionID,
		model.Scopes,
		model.Token,
		model.IssuedAt,
		model.ExpiresAt,
		model.IsRevoked,
	)
	if err != nil {
		return "", err
	}

	return model.RefreshTokenID, nil
}

func (r *oauth2RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*oauth2models.Oauth2RefreshTokenModel, error) {
	const query = `
		SELECT refresh_token_id, user_id, client_id, oidc_session_id, scopes,
			token, issued_at, expires_at, is_revoked, created_at
		FROM oauth2_refresh_tokens
		WHERE token = $1
	`

	var model oauth2models.Oauth2RefreshTokenModel
	row := r.db.QueryRowContext(ctx, query, token)
	err := row.Scan(
		&model.RefreshTokenID,
		&model.UserID,
		&model.ClientID,
		&model.OidcSessionID,
		&model.Scopes,
		&model.Token,
		&model.IssuedAt,
		&model.ExpiresAt,
		&model.IsRevoked,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &model, nil
}

func (r *oauth2RefreshTokenRepository) Revoke(ctx context.Context, refreshTokenID string) error {
	const query = `
		UPDATE oauth2_refresh_tokens
		SET is_revoked = true
		WHERE refresh_token_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, refreshTokenID)
	return err
}

func (r *oauth2RefreshTokenRepository) RevokeByOidcSessionID(ctx context.Context, oidcSessionID string) error {
	const query = `
		UPDATE oauth2_refresh_tokens
		SET is_revoked = true
		WHERE oidc_session_id = $1 AND is_revoked = false
	`

	_, err := r.db.ExecContext(ctx, query, oidcSessionID)
	return err
}

func (r *oauth2RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	const query = `
		UPDATE oauth2_refresh_tokens
		SET is_revoked = true
		WHERE user_id = $1 AND is_revoked = false
	`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
