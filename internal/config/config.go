package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

const directory = "config"

type (
	CustomTime struct {
		Time time.Duration
	}

	// App holds application level configuration.
	App struct {
		Env      string         `json:"env"`
		Debug    string         `json:"debug"`
		HTTP     HTTPConfig     `json:"http"`
		Mysql    SQLConfig      `json:"mysql"`
		Mongo    MongoConfig    `json:"mongo"`
		Log      LogConfig      `json:"log"`
		RabbitMQ RabbitMQConfig `json:"rabbitmq"`
	}

	// MongoConfig holds connection details of mongo server.
	MongoConfig struct {
		AuthSource        string     `json:"authSource"`
		Hosts             []string   `json:"hosts"`
		User              string     `json:"user"`
		Password          string     `json:"password"`
		DB                string     `json:"db"`
		ConnectionTimeout CustomTime `json:"connectionTimeout"`
	}

	// SQLConnConfig holds authentication details for a single sql server.
	SQLConnConfig struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DB       string `json:"db"`
	}

	// SQLConfig holds auth details of all servers.
	SQLConfig struct {
		Servers           []SQLConnConfig `json:"servers"`
		ConnectionTimeout CustomTime      `json:"connectionTimeout"`
	}

	// HTTPConfig holds configuration for HTTP server.
	HTTPConfig struct {
		Host               string     `json:"host"`
		Port               int        `json:"port"`
		ReadTimeout        CustomTime `json:"readTimeout"`
		WriteTimeout       CustomTime `json:"writeTimeout"`
		IdleTimeout        CustomTime `json:"idleTimeout"`
		ShutdownTimeout    CustomTime `json:"shutdownTimeout"`
		MaxHeaderMegabytes int        `json:"maxHeaderMegaBytes"`
	}

	// LogConfig holds configuration for logger.
	LogConfig struct {
		Dir   string `json:"dir"`
		Level int    `json:"level"`
	}

	// RabbitMQConfig holds conf fro amqp server.
	RabbitMQConfig struct {
		Host     string     `json:"host"`
		Port     int        `json:"port"`
		User     string     `json:"user"`
		Password string     `json:"password"`
		Vhost    string     `json:"vhost"`
		Timeout  CustomTime `json:"timeout"`
		Enabled  bool       `json:"enabled"`
	}
)

// Custom unmarshaling for timr in string.
func (ct *CustomTime) UnmarshalJSON(data []byte) (err error) {
	var tmp string

	if err = json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "Custom Time Unmarshal")
	}

	dur, err := time.ParseDuration(tmp)
	if err != nil {
		return errors.Wrap(err, "Custom Time Unmarshal")
	}

	ct.Time = dur

	return errors.Wrap(err, "Custom Time Unmarshal")
}

// Load func loads configuration from *.config.json.
func Load() (*App, error) {
	env := "local"
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	var config *App

	configFile, err := os.Open(filepath.Join(filepath.Clean(directory), filepath.Clean(env)+".config.json"))
	if err != nil {
		return config, errors.Wrap(err, "Config Init")
	}

	defer func() {
		cerr := configFile.Close()
		if cerr == nil {
			err = cerr
		}
	}()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	return config, errors.Wrap(err, "Config Init")
}
