package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type DB struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	DB   string `json:"db"`
	User string `json:"user,omitempty"`
	Pass string `json:"pass,omitempty"`
}

type Logging struct {
	Dir   string `json:"dir"`
	Level int    `json:"level"`
}

type AppConfig struct {
	Env    string `json:"env"`
	Debug  string `json:"debug"`
	Domain string `json:"domian"`
	Port   int    `json:"port"`
	Mysql  DB     `json:"mysql"`
	Mongo  DB     `json:"mongo"`
	Redis  struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		DB   int    `json:"db"`
	} `json:"redis"`
	Log Logging `json:"log"`
}

func Load(env string) (*AppConfig, error) {
	var Config AppConfig

	configFile, err := os.Open(env + ".config.json")
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
