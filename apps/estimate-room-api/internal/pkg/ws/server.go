// Package ws is a websockets implementation
package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type WsServer struct {
	pubClient *redis.Client
	subClient *redis.Client
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewWsServer(pubClient, subClient *redis.Client) *WsServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &WsServer{
		pubClient: pubClient,
		subClient: subClient,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (s *WsServer) Subscribe(channel string, onMessage func([]byte)) {
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

func (s *WsServer) Publish(channel string, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.pubClient.Publish(s.ctx, channel, data).Err()
}

func (s *WsServer) Shutdown() {
	logger.L().Info("Shutting down WebSocket server...")
	s.cancel()

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
