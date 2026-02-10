// Package users provides users endpoints
package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
)

type UsersModule struct {
	Controller UsersController
	Service    UsersService
}

type UsersModuleDeps struct {
	Router      chi.Router
	DB          *pgxpool.Pool
	AuthService auth.AuthService
}

func NewUsersModule(deps UsersModuleDeps) *UsersModule {
	userRepo := usersrepositories.NewUserRepository(deps.DB)
	svc := NewUsersService(deps.AuthService, userRepo)
	ctrl := NewUsersController(svc)

	deps.Router.Route("/users", func(r chi.Router) {
		r.Get("/me", ctrl.GetMe)
	})

	return &UsersModule{
		Controller: ctrl,
		Service:    svc,
	}
}
