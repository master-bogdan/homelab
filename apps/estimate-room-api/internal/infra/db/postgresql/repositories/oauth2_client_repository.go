package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type oauth2ClientRepository struct {
	db *pgxpool.Pool
}

func NewOauth2ClientRepository(db *pgxpool.Pool) Oauth2ClientRepository {
	return &oauth2ClientRepository{db: db}
}

func (r *oauth2ClientRepository) FindByID(clientID string) (*Oauth2ClientModel, error) {
	const query = `
		SELECT client_id, client_secret, redirect_uris, grant_types, response_types, scopes, client_name, client_type, created_at
		FROM oauth2_clients
		WHERE client_id = $1
	`

	var client Oauth2ClientModel
	row := r.db.QueryRow(context.Background(), query, clientID)
	err := row.Scan(
		&client.ClientID,
		&client.ClientSecret,
		&client.RedirectURIs,
		&client.GrantTypes,
		&client.ResponseTypes,
		&client.Scopes,
		&client.ClientName,
		&client.ClientType,
		&client.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return &client, nil
}
