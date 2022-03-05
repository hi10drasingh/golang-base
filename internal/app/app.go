package app

import (
	"flag"

	"github.com/droomlab/drm-coupon/internal/appcontext"
)

func Run(configDir string) error {

	var env string = "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()


	appCtx, err := appcontext.InitilizeAppContext(configDir, env)
	if err != nil {
		return err
	}
	defer appCtx.Close()


	srv := server.NewServer(cfg, handlers.Init(cfg))
}