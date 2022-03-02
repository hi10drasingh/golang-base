package appcontext

import (
	"flag"

	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/droomlab/drm-coupon/pkg/logger"
)

var appContext AppContext
var err error

type AppContext struct {
	Config *config.AppConfig
	Log    logger.Logger
}

func InitilizeAppContext() (*AppContext, error) {
	var env string = "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	appContext.Config, err = config.Load(env)
	if err != nil {
		return &appContext, err
	}

	appContext.Log, err = logger.NewZeroLogger(appContext.Config.Log)

	return &appContext, nil
}
