package oauth2repositories

import (
	"context"
	"database/sql"
	"errors"

	oauth2models "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/models"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type Oauth2ClientRepository interface {
	FindByID(clientID string) (*oauth2models.Oauth2ClientModel, error)
}

type oauth2ClientRepository struct {
	db *bun.DB
}

func NewOauth2ClientRepository(db *bun.DB) *oauth2ClientRepository {
	return &oauth2ClientRepository{db: db}
}

func (r *oauth2ClientRepository) FindByID(clientID string) (*oauth2models.Oauth2ClientModel, error) {
	client := new(oauth2models.Oauth2ClientModel)
	err := r.db.NewSelect().
		Model(client).
		Where("oc.client_id = ?", clientID).
		Limit(1).
		Scan(context.Background())
		// TODO: Refactor this
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrClientNotFound
		}
		return nil, err
	}

	return client, nil
}
