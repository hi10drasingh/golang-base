package handlers

import (
	"net/http"
	"os"

	"github.com/droomlab/drm-coupon/internal/app"
	v1 "github.com/droomlab/drm-coupon/internal/app/handlers/v1"
	"github.com/droomlab/drm-coupon/internal/app/middlewares"
	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/tsenart/nap"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config holds global dependencies for handler
type Config struct {
	Shutdown  chan os.Signal
	Log       drmlog.Logger
	AppConfig *config.App
	SQL       *nap.DB
	NoSQL     *mongo.Client
	RMQ       *drmrmq.RabbitMQ
}

// NewHandlers return new instance of handler
func NewHandlers(conf Config) http.Handler {
	app := app.NewApp(
		conf.Shutdown,
		conf.Log,
		conf.AppConfig,
		middlewares.CORS(),
		middlewares.Logger(conf.Log),
		middlewares.Errors(conf.Log),
		middlewares.Recovery(),
	)

	v1.Routes(app, v1.Config{
		Log:       conf.Log,
		AppConfig: conf.AppConfig,
		SQL:       conf.SQL,
		NoSQL:     conf.NoSQL,
		RMQ:       conf.RMQ,
	})

	return app
}
