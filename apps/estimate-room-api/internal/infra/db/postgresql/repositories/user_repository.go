package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(userID string) (*UserModel, error) {
	const query = `
		SELECT user_id, email, password_hash, github_id, display_name, avatar_url,
			created_at, updated_at, last_login_at, deleted_at
		FROM users
		WHERE user_id = $1
	`

	var user UserModel
	row := r.db.QueryRow(context.Background(), query, userID)
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.GithubID,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*UserModel, error) {
	const query = `
		SELECT user_id, email, password_hash, github_id, display_name, avatar_url,
			created_at, updated_at, last_login_at, deleted_at
		FROM users
		WHERE email = $1
	`

	var user UserModel
	row := r.db.QueryRow(context.Background(), query, email)
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.GithubID,
		&user.DisplayName,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(email, passwordHash string) (string, error) {
	const query = `
		INSERT INTO users (user_id, email, password_hash)
		VALUES ($1, $2, $3)
	`

	userID := uuid.NewString()
	_, err := r.db.Exec(context.Background(), query, userID, email, passwordHash)
	if err != nil {
		return "", err
	}

	return userID, nil
}
