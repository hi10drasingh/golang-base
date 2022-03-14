package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/droomlab/drm-coupon/internal/app/server"
	"github.com/droomlab/drm-coupon/pkg/drmlog"
	"github.com/droomlab/drm-coupon/pkg/drmnosql"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
	"github.com/droomlab/drm-coupon/pkg/drmsql"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

const (
	directory  = "config"
	defaultEnv = "local"
)

type (
	// App holds application level configuration.
	App struct {
		Env      string          `json:"env" validate:"required"`
		Debug    string          `json:"debug" validate:"required,boolean"`
		HTTP     server.Config   `json:"http" validate:"required"`
		Mysql    drmsql.Config   `json:"mysql" validate:"required"`
		Mongo    drmnosql.Config `json:"mongo" validate:"required"`
		Log      drmlog.Config   `json:"log" validate:"required"`
		RabbitMQ drmrmq.Config   `json:"rabbitmq" validate:"required"`
	}
)

func getEnvironment() string {
	env := defaultEnv
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	return env
}

// Load func loads configuration from *.config.json.
func Load() (*App, error) {
	env := getEnvironment()

	var config *App

	configFile, err := os.Open(filepath.Join(filepath.Clean(directory), filepath.Clean(env)+".config.json"))
	if err != nil {
		return nil, errors.Wrap(err, "Config File Open")
	}

	defer func() {
		cerr := configFile.Close()
		if cerr == nil {
			err = cerr
		}
	}()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		return nil, errors.Wrap(err, "Config Json Decode")
	}

	validate := validator.New()

	err = validate.Struct(config)

	if err != nil {
		return nil, errors.Wrap(err, "Config Validation")
	}

	return config, nil
}
