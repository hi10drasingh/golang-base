package appcontext

import (
	"github.com/droomlab/drm-coupon/internal/config"
	"github.com/droomlab/drm-coupon/pkg/logger"
	"github.com/pkg/errors"
)

// AppContext contains all server level dependencies
type AppContext struct {
	Config *config.AppConfig
	Log    logger.Logger
}

// Close used to free memory allocations of dependencies
func (ctx AppContext) Close() {

}

// InitilizeAppContext initializes server level dependencies
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
