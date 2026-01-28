package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type oauth2RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewOauth2RefreshTokenRepository(db *pgxpool.Pool) Oauth2RefreshTokenRepository {
	return &oauth2RefreshTokenRepository{db: db}
}

func (r *oauth2RefreshTokenRepository) Create(model *Oauth2RefreshTokenModel) (string, error) {
	const query = `
		INSERT INTO oauth2_refresh_tokens (
			refresh_token_id, user_id, client_id, oidc_session_id, scopes,
			token, issued_at, expires_at, is_revoked
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if model.RefreshTokenID == "" {
		model.RefreshTokenID = uuid.NewString()
	}

	_, err := r.db.Exec(
		context.Background(),
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

func (r *oauth2RefreshTokenRepository) FindByToken(token string) (*Oauth2RefreshTokenModel, error) {
	const query = `
		SELECT refresh_token_id, user_id, client_id, oidc_session_id, scopes,
			token, issued_at, expires_at, is_revoked, created_at
		FROM oauth2_refresh_tokens
		WHERE token = $1
	`

	var model Oauth2RefreshTokenModel
	row := r.db.QueryRow(context.Background(), query, token)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return &model, nil
}

func (r *oauth2RefreshTokenRepository) Revoke(refreshTokenID string) error {
	const query = `
		UPDATE oauth2_refresh_tokens
		SET is_revoked = true
		WHERE refresh_token_id = $1
	`

	_, err := r.db.Exec(context.Background(), query, refreshTokenID)
	return err
}
