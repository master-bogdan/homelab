package users

import (
	"errors"
	"net/http"

	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	usersmodels "github.com/master-bogdan/estimate-room-api/internal/modules/users/models"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/errors"
)

var ErrUnauthorized = errors.New("unauthorized")

type UsersService interface {
	GetCurrentUser(r *http.Request) (*usersmodels.UserModel, error)
	FindByID(userID string) (*usersmodels.UserModel, error)
	FindByEmail(email string) (*usersmodels.UserModel, error)
	FindByGithubID(githubID string) (*usersmodels.UserModel, error)
	Create(email, passwordHash string) (string, error)
	CreateWithGithub(email *string, githubID, displayName string, avatarURL *string) (string, error)
	UpdateGithubProfile(userID, githubID, displayName string, avatarURL *string, email *string) error
}

type usersService struct {
	authService auth.AuthService
	userRepo    usersrepositories.UserRepository
}

func NewUsersService(authService auth.AuthService, userRepo usersrepositories.UserRepository) UsersService {
	return &usersService{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (s *usersService) GetCurrentUser(r *http.Request) (*usersmodels.UserModel, error) {
	userID, err := s.authService.CheckAuth(r)
	if err != nil {
		if errors.Is(err, auth.ErrMissingToken) {
			return nil, err
		}
		return nil, ErrUnauthorized
	}

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
