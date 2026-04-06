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

type Oauth2AuthCodeRepository interface {
	Create(model *oauth2models.Oauth2AuthCodeModel) error
	FindByCode(code string) (*oauth2models.Oauth2AuthCodeModel, error)
	MarkUsed(authCodeID string) error
	MarkUsedByOidcSessionID(oidcSessionID string) error
	MarkUsedByUserID(userID string) error
}

type oauth2AuthCodeRepository struct {
	db *bun.DB
}

func NewOauth2AuthCodeRepository(db *bun.DB) *oauth2AuthCodeRepository {
	return &oauth2AuthCodeRepository{db: db}
}

func (r *oauth2AuthCodeRepository) Create(model *oauth2models.Oauth2AuthCodeModel) error {
	if model.AuthCodeID == "" {
		model.AuthCodeID = uuid.NewString()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column(
			"auth_code_id",
			"client_id",
			"user_id",
			"oidc_session_id",
			"code",
			"redirect_uri",
			"scopes",
			"code_challenge",
			"code_challenge_method",
			"is_used",
			"expires_at",
		).
		Exec(context.Background())
	return err
}

func (r *oauth2AuthCodeRepository) FindByCode(code string) (*oauth2models.Oauth2AuthCodeModel, error) {
	model := new(oauth2models.Oauth2AuthCodeModel)
	err := r.db.NewSelect().
		Model(model).
		Where("oac.code = ?", code).
		Limit(1).
		Scan(context.Background())
		// TODO: Refactor this
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrAuthCodeNotFound
		}
		return nil, err
	}

	return model, nil
}

func (r *oauth2AuthCodeRepository) MarkUsed(authCodeID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2AuthCodeModel)(nil)).
		Set("is_used = ?", true).
		Where("auth_code_id = ?", authCodeID).
		Exec(context.Background())
	return err
}

func (r *oauth2AuthCodeRepository) MarkUsedByOidcSessionID(oidcSessionID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2AuthCodeModel)(nil)).
		Set("is_used = ?", true).
		Where("oidc_session_id = ?", oidcSessionID).
		Where("is_used = ?", false).
		Exec(context.Background())
	return err
}

func (r *oauth2AuthCodeRepository) MarkUsedByUserID(userID string) error {
	_, err := r.db.NewUpdate().
		Model((*oauth2models.Oauth2AuthCodeModel)(nil)).
		Set("is_used = ?", true).
		Where("user_id = ?", userID).
		Where("is_used = ?", false).
		Exec(context.Background())
	return err
}
