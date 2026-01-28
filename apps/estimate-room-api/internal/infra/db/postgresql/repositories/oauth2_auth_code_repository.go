package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type oauth2AuthCodeRepository struct {
	db *pgxpool.Pool
}

func NewOauth2AuthCodeRepository(db *pgxpool.Pool) Oauth2AuthCodeRepository {
	return &oauth2AuthCodeRepository{db: db}
}

func (r *oauth2AuthCodeRepository) Create(model *Oauth2AuthCodeModel) error {
	const query = `
		INSERT INTO oauth2_auth_codes (
			auth_code_id, client_id, user_id, oidc_session_id, code,
			redirect_uri, scopes, code_challenge, code_challenge_method, is_used, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	if model.AuthCodeID == "" {
		model.AuthCodeID = uuid.NewString()
	}

	_, err := r.db.Exec(
		context.Background(),
		query,
		model.AuthCodeID,
		model.ClientID,
		model.UserID,
		model.OidcSessionID,
		model.Code,
		model.RedirectURI,
		model.Scopes,
		model.CodeChallenge,
		model.CodeChallengeMethod,
		model.IsUsed,
		model.ExpiresAt,
	)
	return err
}

func (r *oauth2AuthCodeRepository) FindByCode(code string) (*Oauth2AuthCodeModel, error) {
	const query = `
		SELECT auth_code_id, client_id, user_id, oidc_session_id, code,
			redirect_uri, scopes, code_challenge, code_challenge_method, is_used, expires_at, created_at
		FROM oauth2_auth_codes
		WHERE code = $1
	`

	var model Oauth2AuthCodeModel
	row := r.db.QueryRow(context.Background(), query, code)
	err := row.Scan(
		&model.AuthCodeID,
		&model.ClientID,
		&model.UserID,
		&model.OidcSessionID,
		&model.Code,
		&model.RedirectURI,
		&model.Scopes,
		&model.CodeChallenge,
		&model.CodeChallengeMethod,
		&model.IsUsed,
		&model.ExpiresAt,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAuthCodeNotFound
		}
		return nil, err
	}

	return &model, nil
}

func (r *oauth2AuthCodeRepository) MarkUsed(authCodeID string) error {
	const query = `
		UPDATE oauth2_auth_codes
		SET is_used = true
		WHERE auth_code_id = $1
	`

	_, err := r.db.Exec(context.Background(), query, authCodeID)
	return err
}
