package appcontext

import (
	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/logger"
	"github.com/pkg/errors"
)

type AppContext struct {
	Config *config.AppConfig
	Log    logger.Logger
}

func (ctx AppContext) Close() {

}

func InitilizeAppContext(configDir string, env string) (*AppContext, error) {
	conf, err := config.Load(configDir, env)
	if err != nil {
		return nil, errors.Wrap(err, "Config Initialize")
	}

	log, err := logger.NewZeroLogger(conf.Log)
	if err != nil {
		return nil, errors.Wrap(err, "Log Initialize")
	}

	return &AppContext{conf, log}, nil
}
