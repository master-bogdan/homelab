// Package users provides users endpoints
package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/internal/modules/auth"
)

type UsersModule struct {
	Controller UsersController
	Service    UsersService
}

type UsersModuleDeps struct {
	Router      chi.Router
	AuthService auth.AuthService
	UserRepo    UserRepository
}

func NewUsersModule(deps UsersModuleDeps) *UsersModule {
	svc := NewUsersService(deps.AuthService, deps.UserRepo)
	ctrl := NewUsersController(svc)

	deps.Router.Route("/users", func(r chi.Router) {
		r.Get("/me", ctrl.GetMe)
	})

	return &UsersModule{
		Controller: ctrl,
		Service:    svc,
	}
}
