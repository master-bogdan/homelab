package usersrepositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

type UserRepository interface {
	FindByID(userID string) (*usersmodels.UserModel, error)
	FindByEmail(email string) (*usersmodels.UserModel, error)
	FindByGithubID(githubID string) (*usersmodels.UserModel, error)
	Create(email, passwordHash string) (string, error)
	CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error)
	UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(userID string) (*usersmodels.UserModel, error) {
	const query = `
		SELECT user_id, email, password_hash, github_id, display_name, avatar_url,
			created_at, updated_at, last_login_at, deleted_at
		FROM users
		WHERE user_id = $1
	`

	var user usersmodels.UserModel
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

	if err == nil {
		return &user, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrUserNotFound
	}

	return nil, err
}

func (r *userRepository) FindByEmail(email string) (*usersmodels.UserModel, error) {
	const query = `
		SELECT user_id, email, password_hash, github_id, display_name, avatar_url,
			created_at, updated_at, last_login_at, deleted_at
		FROM users
		WHERE email = $1
	`

	var user usersmodels.UserModel
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

	if err == nil {
		return &user, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrUserNotFound
	}

	return nil, err
}

func (r *userRepository) FindByGithubID(githubID string) (*usersmodels.UserModel, error) {
	const query = `
		SELECT user_id, email, password_hash, github_id, display_name, avatar_url,
			created_at, updated_at, last_login_at, deleted_at
		FROM users
		WHERE github_id = $1
	`

	var user usersmodels.UserModel
	row := r.db.QueryRow(context.Background(), query, githubID)
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

	if err == nil {
		return &user, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.ErrUserNotFound
	}

	return nil, err
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

func (r *userRepository) CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error) {
	const query = `
		INSERT INTO users (user_id, email, github_id, display_name, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
	`

	userID := uuid.NewString()
	_, err := r.db.Exec(context.Background(), query, userID, email, githubID, displayName, avatarURL)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *userRepository) UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error {
	const query = `
		UPDATE users
		SET github_id = $2,
			display_name = $3,
			avatar_url = $4,
			email = COALESCE($5, email),
			updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.Exec(context.Background(), query, userID, githubID, displayName, avatarURL, email)
	return err
}
