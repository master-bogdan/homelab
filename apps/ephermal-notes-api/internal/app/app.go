package app

import (
	"log/slog"
	"net/http"

	"github.com/master-bogdan/ephermal-notes/internal/infra/metrics"
	"github.com/master-bogdan/ephermal-notes/internal/modules/health"
	"github.com/master-bogdan/ephermal-notes/internal/modules/notes"
	_ "github.com/master-bogdan/ephermal-notes/pkg/docs"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	Router *http.ServeMux
	Client *redis.Client
	Logger *slog.Logger
}

func Init(app App) {
	app.Router.Handle("/api/v1/swagger/", httpSwagger.WrapHandler)
	app.Router.Handle("/api/v1/metrics", metrics.Handler())

	health.RouterNew(app.Router, app.Client, app.Logger)
	notes.RouterNew(app.Router, app.Client, app.Logger)
}
