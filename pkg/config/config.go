package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type (
	customTime time.Duration

	// AppConfig holds application level configuration
	AppConfig struct {
		Env   string     `json:"env"`
		Debug string     `json:"debug"`
		HTTP  HTTPConfig `json:"http"`
		Mysql DBConfig   `json:"mysql"`
		Mongo DBConfig   `json:"mongo"`
		Log   LogConfig  `json:"log"`
	}

	// DBConfig holds configuration for Database server
	DBConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		DB   string `json:"db"`
		User string `json:"user,omitempty"`
		Pass string `json:"pass,omitempty"`
	}

	// HTTPConfig holds configuration for HTTP server
	HTTPConfig struct {
		Host               string     `json:"host"`
		Port               int        `json:"port"`
		ReadTimeout        customTime `json:"readTimeout"`
		WriteTimeout       customTime `json:"writeTimeout"`
		IdleTimeout        customTime `json:"idleTimeout"`
		MaxHeaderMegabytes int        `json:"maxHeaderMegaBytes"`
	}

	// LogConfig holds configuration for logger
	LogConfig struct {
		Dir   string `json:"dir"`
		Level int    `json:"level"`
	}
)

// Custom unmarshaling for timr in string
func (c *customTime) UnmarshalJSON(data []byte) (err error) {
	var tmp string

	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	time, err := time.ParseDuration(tmp)

	*c = customTime(time)

	return err
}

// Load func loads configuration from *.config.json
func Load(dir string, env string) (*AppConfig, error) {
	var config *AppConfig

	configFile, err := os.Open(filepath.Join(filepath.Clean(dir), filepath.Clean(env)+".config.json"))
	if err != nil {
		return config, errors.Wrap(err, "Config File Open")
	}

	defer func() {
		cerr := configFile.Close()
		if err == nil {
			err = cerr
		}
	}()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		return config, errors.Wrap(err, "Config File Decode")
	}

	return config, err
}
