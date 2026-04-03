package authrepositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	authmodels "github.com/master-bogdan/estimate-room-api/internal/modules/auth/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, model *authmodels.PasswordResetTokenModel) (string, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*authmodels.PasswordResetTokenModel, error)
	MarkUsed(ctx context.Context, passwordResetTokenID string) error
}

type passwordResetTokenRepository struct {
	db *bun.DB
}

func NewPasswordResetTokenRepository(db *bun.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, model *authmodels.PasswordResetTokenModel) (string, error) {
	if model.PasswordResetTokenID == "" {
		model.PasswordResetTokenID = uuid.NewString()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("password_reset_token_id", "user_id", "token_hash", "expires_at", "used_at").
		Exec(ctx)
	if err != nil {
		return "", err
	}

	return model.PasswordResetTokenID, nil
}

func (r *passwordResetTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*authmodels.PasswordResetTokenModel, error) {
	model := new(authmodels.PasswordResetTokenModel)
	err := r.db.NewSelect().
		Model(model).
		Where("prt.token_hash = ?", tokenHash).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrPasswordResetTokenNotFound
		}
		return nil, err
	}

	return model, nil
}

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, passwordResetTokenID string) error {
	_, err := r.db.NewUpdate().
		Model((*authmodels.PasswordResetTokenModel)(nil)).
		Set("used_at = COALESCE(used_at, NOW())").
		Where("password_reset_token_id = ?", passwordResetTokenID).
		Exec(ctx)

	return err
}
