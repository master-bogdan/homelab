package authrepositories

import (
	"context"
	"errors"

	"database/sql"
	"github.com/google/uuid"
	authmodels "github.com/master-bogdan/estimate-room-api/internal/modules/auth/models"
	"github.com/uptrace/bun"
)

type OidcSessionRepository interface {
	Create(model *authmodels.OidcSessionModel) (string, error)
	FindByID(sessionID string) (*authmodels.OidcSessionModel, error)
}

type oauth2OidcSessionRepository struct {
	db *bun.DB
}

func NewOauth2OidcSessionRepository(db *bun.DB) *oauth2OidcSessionRepository {
	return &oauth2OidcSessionRepository{db: db}
}

func (r *oauth2OidcSessionRepository) Create(model *authmodels.OidcSessionModel) (string, error) {
	const query = `
		INSERT INTO oauth2_oidc_sessions (oidc_session_id, user_id, client_id, nonce)
		VALUES ($1, $2, $3, $4)
	`

	if model.OidcSessionID == "" {
		model.OidcSessionID = uuid.NewString()
	}

	_, err := r.db.ExecContext(
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

func (r *oauth2OidcSessionRepository) FindByID(sessionID string) (*authmodels.OidcSessionModel, error) {
	const query = `
		SELECT oidc_session_id, user_id, client_id, nonce, created_at
		FROM oauth2_oidc_sessions
		WHERE oidc_session_id = $1
	`

	var model authmodels.OidcSessionModel
	row := r.db.QueryRowContext(context.Background(), query, sessionID)
	err := row.Scan(&model.OidcSessionID, &model.UserID, &model.ClientID, &model.Nonce, &model.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOidcSessionNotFound
		}
		return nil, err
	}

	return &model, nil
}
