package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
)

const (
	directory  = "config"
	defaultEnv = "local"
)

type (
	CustomTime struct {
		Time time.Duration `validate:"required"`
	}

	// App holds application level configuration.
	App struct {
		Env      string         `json:"env" validate:"required"`
		Debug    string         `json:"debug" validate:"required,boolean"`
		HTTP     HTTPConfig     `json:"http" validate:"required"`
		Mysql    SQLConfig      `json:"mysql" validate:"required"`
		Mongo    MongoConfig    `json:"mongo" validate:"required"`
		Log      LogConfig      `json:"log" validate:"required"`
		RabbitMQ RabbitMQConfig `json:"rabbitmq" validate:"required"`
	}

	// MongoConfig holds connection details of mongo server.
	MongoConfig struct {
		// AuthSource        string     `json:"authSource" validate:"required"`
		Hosts             []string   `json:"hosts" validate:"required"`
		User              string     `json:"user" validate:"required"`
		Password          string     `json:"password" validate:"required"`
		DB                string     `json:"db" validate:"required"`
		ConnectionTimeout CustomTime `json:"connectionTimeout" validate:"required"`
	}

	// SQLConnConfig holds authentication details for a single sql server.
	SQLConnConfig struct {
		Host     string `json:"host" validate:"required"`
		Port     int    `json:"port" validate:"required,number"`
		User     string `json:"user" validate:"required"`
		Password string `json:"password" validate:"required"`
		DB       string `json:"db" validate:"required"`
	}

	// SQLConfig holds auth details of all servers.
	SQLConfig struct {
		Servers           []SQLConnConfig `json:"servers" validate:"required"`
		ConnectionTimeout CustomTime      `json:"connectionTimeout" validate:"required"`
	}

	// HTTPConfig holds configuration for HTTP server.
	HTTPConfig struct {
		// Host               string     `json:"host" validate:"required"`
		Port               int        `json:"port" validate:"required,number"`
		ReadTimeout        CustomTime `json:"readTimeout" validate:"required"`
		WriteTimeout       CustomTime `json:"writeTimeout" validate:"required"`
		IdleTimeout        CustomTime `json:"idleTimeout" validate:"required"`
		ShutdownTimeout    CustomTime `json:"shutdownTimeout" validate:"required"`
		MaxHeaderMegabytes int        `json:"maxHeaderMegaBytes" validate:"required,number"`
	}

	// LogConfig holds configuration for logger.
	LogConfig struct {
		Dir   string `json:"dir" validate:"required"`
		Level int    `json:"level" validate:"number"`
	}

	// RabbitMQConfig holds conf fro amqp server.
	RabbitMQConfig struct {
		Host     string     `json:"host" validate:"required"`
		Port     int        `json:"port" validate:"required,number"`
		User     string     `json:"user" validate:"required"`
		Password string     `json:"password" validate:"required"`
		Vhost    string     `json:"vhost" validate:"required"`
		Timeout  CustomTime `json:"timeout" validate:"required"`
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
