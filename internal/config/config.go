package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/droomlab/drm-coupon/pkg/logger"
	"github.com/pkg/errors"
)

type(
	AppConfig struct {
		Env    string `json:"env"`
		Debug  string `json:"debug"`
		HTTP HTTPConfig `json:"http"`
		Mysql  DBConfig     `json:"mysql"`
		Mongo  DBConfig     `json:"mongo"`
		Log logger.LogConfig `json:"log"`
	}

	DBConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		DB   string `json:"db"`
		User string `json:"user,omitempty"`
		Pass string `json:"pass,omitempty"`
	}

	HTTPConfig struct {
		Host string `json:"host"`
		Port int `json:"port"`
		ReadTimeout time.Duration `json:"readTimeout"`
		WriteTimeout time.Duration `json:"writeTimeout"`
		IdleTimeout time.Duration `json:"idleTimeout"`
		MaxHeaderMegabytes int `json:"maxHeaderMegaBytes"`
	}
)

func Load(dir string, env string) (*AppConfig, error) {
	var Config AppConfig

	configFile, err := os.Open(dir + "/" + env + ".config.json")
	if err != nil {
		return nil, errors.Wrap(err, "Config File Open")
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)

	if err != nil {
		return nil, errors.Wrap(err, "Config File Decode")
	}

	return &Config, nil
}
