package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/droomlab/drm-coupon/pkg/logger"
	"github.com/pkg/errors"
)

type (
	CustomTime time.Duration
	AppConfig  struct {
		Env   string           `json:"env"`
		Debug string           `json:"debug"`
		HTTP  HTTPConfig       `json:"http"`
		Mysql DBConfig         `json:"mysql"`
		Mongo DBConfig         `json:"mongo"`
		Log   logger.LogConfig `json:"log"`
	}

	DBConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		DB   string `json:"db"`
		User string `json:"user,omitempty"`
		Pass string `json:"pass,omitempty"`
	}

	HTTPConfig struct {
		Host               string     `json:"host"`
		Port               int        `json:"port"`
		ReadTimeout        CustomTime `json:"readTimeout"`
		WriteTimeout       CustomTime `json:"writeTimeout"`
		IdleTimeout        CustomTime `json:"idleTimeout"`
		MaxHeaderMegabytes int        `json:"maxHeaderMegaBytes"`
	}
)

func (c *CustomTime) UnmarshalJSON(data []byte) (err error) {
	var tmp string

	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	time, err := time.ParseDuration(tmp)

	*c = CustomTime(time)

	return err
}

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
