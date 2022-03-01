package appcontext

import (
	"flag"

	"github.com/droomlab/drm-coupon/pkg/config"
)

var appContext AppContext
var err error

type AppContext struct {
	Config *config.AppConfig
}

func InitilizeAppContext() (*AppContext, error) {
	var env string = "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()


	appContext.Config, err = config.Load(env)
	if err != nil {
		return &appContext,err
	}

	return &appContext,nil
}
