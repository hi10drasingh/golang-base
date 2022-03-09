package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

const directory string = "config"

type (
	customTime time.Duration

	// App holds application level configuration
	App struct {
		Env   string      `json:"env"`
		Debug string      `json:"debug"`
		HTTP  HTTPConfig  `json:"http"`
		Mysql SQLConfig   `json:"mysql"`
		Mongo MongoConfig `json:"mongo"`
		Log   LogConfig   `json:"log"`
	}

	// MongoConfig holds connection details of mongo server
	MongoConfig struct {
		Code          string     `json:"code"`
		AuthMechanism string     `json:"authMechanism"`
		Host          string     `json:"host"`
		Port          int        `json:"port"`
		User          string     `json:"user"`
		Password      string     `json:"password"`
		DB            string     `json:"db"`
		Timeout       customTime `json:"timeout"`
	}

	// SQLConnConfig holds authentication details for a single sql server
	SQLConnConfig struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		DB       string `json:"db"`
	}

	// SQLConfig holds auth details of all servers
	SQLConfig []SQLConnConfig

	// HTTPConfig holds configuration for HTTP server
	HTTPConfig struct {
		Host               string     `json:"host"`
		Port               int        `json:"port"`
		ReadTimeout        customTime `json:"readTimeout"`
		WriteTimeout       customTime `json:"writeTimeout"`
		IdleTimeout        customTime `json:"idleTimeout"`
		ShutdownTimeout    customTime `json:"shutdownTimeout"`
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
func Load() (*App, error) {
	var env string = "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	var config *App

	configFile, err := os.Open(filepath.Join(filepath.Clean(directory), filepath.Clean(env)+".config.json"))
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
