package appcontext

import (
	"flag"

	"github.com/droomlab/drm-coupon/pkg/config"
	"github.com/droomlab/drm-coupon/pkg/logger"
	"github.com/pkg/errors"
)

type AppContext struct {
	Config *config.AppConfig
	Log    logger.Logger
}

func InitilizeAppContext() (*AppContext, error) {
	var env string = "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	conf, err := config.Load(env)
	if err != nil {
		return nil, errors.Wrap(err, "Config Initialize")
	}

	log, err := logger.NewZeroLogger(conf.Log)
	if err != nil {
		return nil, errors.Wrap(err, "Log Initialize")
	}

	return &AppContext{conf, log}, nil
}
