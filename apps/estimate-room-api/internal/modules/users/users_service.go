package users

import (
	"errors"
	"net/http"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
)

var ErrUnauthorized = errors.New("unauthorized")

type UsersService interface {
	GetCurrentUser(r *http.Request) (*models.UserModel, error)
}

type usersService struct {
	authService auth.AuthService
	userRepo    UserRepository
}

func NewUsersService(authService auth.AuthService, userRepo UserRepository) UsersService {
	return &usersService{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (s *usersService) GetCurrentUser(r *http.Request) (*models.UserModel, error) {
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
		return nil, repositories.ErrUserNotFound
	}

	return user, nil
}
