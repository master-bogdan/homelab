package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type oauth2OidcSessionRepository struct {
	db *pgxpool.Pool
}

func NewOauth2OidcSessionRepository(db *pgxpool.Pool) Oauth2OidcSessionRepository {
	return &oauth2OidcSessionRepository{db: db}
}

func (r *oauth2OidcSessionRepository) Create(model *OidcSessionModel) (string, error) {
	const query = `
		INSERT INTO oauth2_oidc_sessions (oidc_session_id, user_id, client_id, nonce)
		VALUES ($1, $2, $3, $4)
	`

	if model.OidcSessionID == "" {
		model.OidcSessionID = uuid.NewString()
	}

	_, err := r.db.Exec(
		context.Background(),
		query,
		model.OidcSessionID,
		model.UserID,
		model.ClientID,
		model.Nonce,
	)
	if err != nil {
		return "", err
	}

	return model.OidcSessionID, nil
}

func (r *oauth2OidcSessionRepository) FindByID(sessionID string) (*OidcSessionModel, error) {
	const query = `
		SELECT oidc_session_id, user_id, client_id, nonce, created_at
		FROM oauth2_oidc_sessions
		WHERE oidc_session_id = $1
	`

	var model OidcSessionModel
	row := r.db.QueryRow(context.Background(), query, sessionID)
	err := row.Scan(&model.OidcSessionID, &model.UserID, &model.ClientID, &model.Nonce, &model.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOidcSessionNotFound
		}
		return nil, err
	}

	return &model, nil
}
