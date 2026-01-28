package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type oauth2UserRepository struct {
	db *pgxpool.Pool
}

func NewOauth2UserRepository(db *pgxpool.Pool) Oauth2UserRepository {
	return &oauth2UserRepository{db: db}
}

func (r *oauth2UserRepository) FindByEmail(email string) (*Oauth2UserModel, error) {
	const query = `
		SELECT user_id, email, password_hash, created_at
		FROM oauth2_users
		WHERE email = $1
	`

	var user Oauth2UserModel
	row := r.db.QueryRow(context.Background(), query, email)
	err := row.Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *oauth2UserRepository) Create(email, passwordHash string) (string, error) {
	const query = `
		INSERT INTO oauth2_users (user_id, email, password_hash)
		VALUES ($1, $2, $3)
	`

	userID := uuid.NewString()
	_, err := r.db.Exec(context.Background(), query, userID, email, passwordHash)
	if err != nil {
		return "", err
	}

	return userID, nil
}
