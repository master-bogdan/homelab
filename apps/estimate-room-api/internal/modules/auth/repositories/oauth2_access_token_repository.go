package authrepositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	authmodels "github.com/master-bogdan/estimate-room-api/internal/modules/auth/models"
)

type AccessTokenRepository interface {
	Create(ctx context.Context, model *authmodels.Oauth2AccessTokenModel) error
	FindByToken(ctx context.Context, token string) (*authmodels.Oauth2AccessTokenModel, error)
	Revoke(ctx context.Context, accessTokenID string) error
}

type oauth2AccessTokenRepository struct {
	db *pgxpool.Pool
}

func NewOauth2AccessTokenRepository(db *pgxpool.Pool) *oauth2AccessTokenRepository {
	return &oauth2AccessTokenRepository{db: db}
}

func (r *oauth2AccessTokenRepository) Create(ctx context.Context, model *authmodels.Oauth2AccessTokenModel) error {
	const query = `
		INSERT INTO oauth2_access_tokens (
			access_token_id, user_id, client_id, oidc_session_id, refresh_token_id,
			scopes, token, issued_at, expires_at, issuer, is_revoked
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	if model.AccessTokenID == "" {
		model.AccessTokenID = uuid.NewString()
	}

	_, err := r.db.Exec(
		ctx,
		query,
		model.AccessTokenID,
		model.UserID,
		model.ClientID,
		model.OidcSessionID,
		model.RefreshTokenID,
		model.Scopes,
		model.Token,
		model.IssuedAt,
		model.ExpiresAt,
		model.Issuer,
		model.IsRevoked,
	)
	return err
}

func (r *oauth2AccessTokenRepository) FindByToken(ctx context.Context, token string) (*authmodels.Oauth2AccessTokenModel, error) {
	const query = `
		SELECT access_token_id, user_id, client_id, oidc_session_id, refresh_token_id,
			scopes, token, issued_at, expires_at, issuer, is_revoked, created_at
		FROM oauth2_access_tokens
		WHERE token = $1
	`

	var model authmodels.Oauth2AccessTokenModel
	row := r.db.QueryRow(ctx, query, token)
	err := row.Scan(
		&model.AccessTokenID,
		&model.UserID,
		&model.ClientID,
		&model.OidcSessionID,
		&model.RefreshTokenID,
		&model.Scopes,
		&model.Token,
		&model.IssuedAt,
		&model.ExpiresAt,
		&model.Issuer,
		&model.IsRevoked,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccessTokenNotFound
		}
		return nil, err
	}

	return &model, nil
}

func (r *oauth2AccessTokenRepository) Revoke(ctx context.Context, accessTokenID string) error {
	const query = `
		UPDATE oauth2_access_tokens
		SET is_revoked = true
		WHERE access_token_id = $1
	`

	_, err := r.db.Exec(ctx, query, accessTokenID)
	return err
}
