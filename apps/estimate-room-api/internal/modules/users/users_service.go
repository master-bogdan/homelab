package users

import (
	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

type UsersService interface {
	GetCurrentUser(userID string) (*usersmodels.UserModel, error)
	FindByID(userID string) (*usersmodels.UserModel, error)
	FindByEmail(email string) (*usersmodels.UserModel, error)
	FindByGithubID(githubID string) (*usersmodels.UserModel, error)
	Create(email, passwordHash string) (string, error)
	CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error)
	UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error
}

type usersService struct {
	userRepo usersrepositories.UserRepository
}

func NewUsersService(userRepo usersrepositories.UserRepository) UsersService {
	return &usersService{
		userRepo: userRepo,
	}
}

func (s *usersService) GetCurrentUser(userID string) (*usersmodels.UserModel, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt != nil {
		return nil, apperrors.ErrUserNotFound
	}

	return user, nil
}

func (s *usersService) FindByID(userID string) (*usersmodels.UserModel, error) {
	return s.userRepo.FindByID(userID)
}

func (s *usersService) FindByEmail(email string) (*usersmodels.UserModel, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *usersService) FindByGithubID(githubID string) (*usersmodels.UserModel, error) {
	return s.userRepo.FindByGithubID(githubID)
}

func (s *usersService) Create(email, passwordHash string) (string, error) {
	return s.userRepo.Create(email, passwordHash)
}

func (s *usersService) CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error) {
	return s.userRepo.CreateWithGithub(email, githubID, displayName, avatarURL)
}

func (s *usersService) UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error {
	return s.userRepo.UpdateGithubProfile(userID, githubID, displayName, avatarURL, email)
}
