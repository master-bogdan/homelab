package app

import (
	"context"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/config"
)

func TestSetupApp_RequiresFrontendBaseURL(t *testing.T) {
	deps := AppDeps{
		Cfg:    &config.Config{},
		Router: chi.NewRouter(),
	}

	err := deps.SetupApp(context.Background())
	if err == nil {
		t.Fatalf("expected missing FRONTEND_BASE_URL to return an error")
	}
	if !strings.Contains(err.Error(), "FRONTEND_BASE_URL") {
		t.Fatalf("expected FRONTEND_BASE_URL error, got %v", err)
	}
}
