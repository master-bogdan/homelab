package usersrepositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
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
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(userID string) (*usersmodels.UserModel, error) {
	user := new(usersmodels.UserModel)
	err := r.db.NewSelect().
		Model(user).
		Where("u.user_id = ?", userID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*usersmodels.UserModel, error) {
	user := new(usersmodels.UserModel)
	err := r.db.NewSelect().
		Model(user).
		Where("u.email = ?", email).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByGithubID(githubID string) (*usersmodels.UserModel, error) {
	user := new(usersmodels.UserModel)
	err := r.db.NewSelect().
		Model(user).
		Where("u.github_id = ?", githubID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(email, passwordHash string) (string, error) {
	user := &usersmodels.UserModel{
		UserID: uuid.NewString(),
		Email:  &email,
	}
	if passwordHash != "" {
		user.PasswordHash = &passwordHash
	}

	_, err := r.db.NewInsert().
		Model(user).
		Column("user_id", "email", "password_hash").
		Exec(context.Background())
	if err != nil {
		return "", err
	}

	return user.UserID, nil
}

func (r *userRepository) CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error) {
	user := &usersmodels.UserModel{
		UserID:      uuid.NewString(),
		Email:       email,
		GithubID:    &githubID,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}

	_, err := r.db.NewInsert().
		Model(user).
		Column("user_id", "email", "github_id", "display_name", "avatar_url").
		Exec(context.Background())
	if err != nil {
		return "", err
	}

	return user.UserID, nil
}

func (r *userRepository) UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error {
	user := &usersmodels.UserModel{UserID: userID}
	_, err := r.db.NewUpdate().
		Model(user).
		Set("github_id = ?", githubID).
		Set("display_name = ?", displayName).
		Set("avatar_url = ?", avatarURL).
		Set("email = COALESCE(?, email)", email).
		Set("updated_at = NOW()").
		Where("user_id = ?", userID).
		Exec(context.Background())

	return err
}
