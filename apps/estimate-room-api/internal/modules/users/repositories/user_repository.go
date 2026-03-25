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
	HasSoftDeletedEmail(email string) (bool, error)
	HasSoftDeletedGithubID(githubID string) (bool, error)
	Create(email, passwordHash, displayName string, organization, occupation *string) (string, error)
	CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error)
	UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error
	UpdateDisplayName(userID, displayName string) error
	UpdatePasswordHash(userID, passwordHash string) error
	UpdateLastLoginAt(userID string) error
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
		Where("u.deleted_at IS NULL").
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
		Where("u.deleted_at IS NULL").
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

func (r *userRepository) HasSoftDeletedEmail(email string) (bool, error) {
	return r.hasSoftDeletedUser("u.email = ?", email)
}

func (r *userRepository) HasSoftDeletedGithubID(githubID string) (bool, error) {
	return r.hasSoftDeletedUser("u.github_id = ?", githubID)
}

func (r *userRepository) hasSoftDeletedUser(whereClause string, value string) (bool, error) {
	return r.db.NewSelect().
		Model((*usersmodels.UserModel)(nil)).
		Where(whereClause, value).
		Where("u.deleted_at IS NOT NULL").
		Exists(context.Background())
}

func (r *userRepository) Create(email, passwordHash, displayName string, organization, occupation *string) (string, error) {
	user := &usersmodels.UserModel{
		UserID:       uuid.NewString(),
		Email:        &email,
		DisplayName:  displayName,
		Organization: organization,
		Occupation:   occupation,
	}
	if passwordHash != "" {
		user.PasswordHash = &passwordHash
	}

	_, err := r.db.NewInsert().
		Model(user).
		Column("user_id", "email", "password_hash", "display_name", "organization", "occupation").
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

func (r *userRepository) UpdateDisplayName(userID, displayName string) error {
	_, err := r.db.NewUpdate().
		Model((*usersmodels.UserModel)(nil)).
		Set("display_name = ?", displayName).
		Set("updated_at = NOW()").
		Where("user_id = ?", userID).
		Exec(context.Background())

	return err
}

func (r *userRepository) UpdatePasswordHash(userID, passwordHash string) error {
	_, err := r.db.NewUpdate().
		Model((*usersmodels.UserModel)(nil)).
		Set("password_hash = ?", passwordHash).
		Set("updated_at = NOW()").
		Where("user_id = ?", userID).
		Exec(context.Background())

	return err
}

func (r *userRepository) UpdateLastLoginAt(userID string) error {
	_, err := r.db.NewUpdate().
		Model((*usersmodels.UserModel)(nil)).
		Set("last_login_at = NOW()").
		Set("updated_at = NOW()").
		Where("user_id = ?", userID).
		Exec(context.Background())

	return err
}
