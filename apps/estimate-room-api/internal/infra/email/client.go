package email

import (
	"context"
	"errors"
)

var ErrNotConfigured = errors.New("email client is not configured")

type Message struct {
	To       []string
	Subject  string
	TextBody string
}

type Client interface {
	Send(ctx context.Context, msg Message) error
}

type noopClient struct{}

func NewNoopClient() Client {
	return noopClient{}
}

func (noopClient) Send(_ context.Context, _ Message) error {
	return ErrNotConfigured
}
