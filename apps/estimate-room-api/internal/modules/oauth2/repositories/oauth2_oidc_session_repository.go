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

type OidcSessionRepository interface {
	Create(model *oauth2models.OidcSessionModel) (string, error)
	FindByID(sessionID string) (*oauth2models.OidcSessionModel, error)
	Revoke(sessionID string) error
	RevokeByUserID(userID string) error
}

type oauth2OidcSessionRepository struct {
	db *bun.DB
}

func NewOauth2OidcSessionRepository(db *bun.DB) *oauth2OidcSessionRepository {
	return &oauth2OidcSessionRepository{db: db}
}

func (r *oauth2OidcSessionRepository) Create(model *oauth2models.OidcSessionModel) (string, error) {
	if model.OidcSessionID == "" {
		model.OidcSessionID = uuid.NewString()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("oidc_session_id", "user_id", "client_id", "nonce", "revoked_at").
		Exec(context.Background())
	if err != nil {
		return "", err
	}

	return model.OidcSessionID, nil
}

func (r *oauth2OidcSessionRepository) FindByID(sessionID string) (*oauth2models.OidcSessionModel, error) {
	model := new(oauth2models.OidcSessionModel)
	err := r.db.NewSelect().
		Model(model).
		Where("os.oidc_session_id = ?", sessionID).
		Where("os.revoked_at IS NULL").
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrOidcSessionNotFound
		}
		return nil, err
	}

	return model, nil
}

func (r *oauth2OidcSessionRepository) Revoke(sessionID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.OidcSessionModel)(nil)).
		Set("revoked_at = COALESCE(revoked_at, NOW())").
		Where("oidc_session_id = ?", sessionID).
		Exec(context.Background())

	return err
}

func (r *oauth2OidcSessionRepository) RevokeByUserID(userID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.OidcSessionModel)(nil)).
		Set("revoked_at = COALESCE(revoked_at, NOW())").
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Exec(context.Background())

	return err
}
