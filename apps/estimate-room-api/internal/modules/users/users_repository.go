package users

import "github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/models"

type UserRepository interface {
	FindByID(userID string) (*models.UserModel, error)
	FindByEmail(email string) (*models.UserModel, error)
	Create(email, passwordHash string) (string, error)
}
