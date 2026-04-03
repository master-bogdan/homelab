// Package wsserver provides websocket server infrastructure.
package wsserver

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	pubClient *redis.Client
	subClient *redis.Client
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

type ServerDeps struct {
	PubClient *redis.Client
	SubClient *redis.Client
}

func NewServer(deps ServerDeps) (*Server, error) {
	if deps.PubClient == nil {
		return nil, errors.New("redis pub client is required")
	}
	if deps.SubClient == nil {
		return nil, errors.New("redis sub client is required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		pubClient: deps.PubClient,
		subClient: deps.SubClient,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (s *Server) Subscribe(channel string, onMessage func([]byte)) {
	if s == nil || s.subClient == nil {
		return
	}
	if onMessage == nil {
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		pubsub := s.subClient.Subscribe(s.ctx, channel)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case <-s.ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				if msg != nil {
					onMessage([]byte(msg.Payload))
				}
			}
		}
	}()
}

func (s *Server) Publish(channel string, message any) error {
	if s == nil || s.pubClient == nil {
		return errors.New("redis pub client is not initialized")
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.pubClient.Publish(s.ctx, channel, data).Err()
}

func (s *Server) Shutdown() {
	if s == nil {
		return
	}

	logger.L().Info("Shutting down WebSocket server...")
	if s.cancel != nil {
		s.cancel()
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.L().Info("WebSocket server shut down gracefully")
	case <-time.After(5 * time.Second):
		logger.L().Warn("WebSocket server shutdown timeout")
	}
}
