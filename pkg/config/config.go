package config

import (
	"encoding/json"
	"os"
)

type DB struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	DB   string `json:"db"`
	User string `json:"user,omitempty"`
	Pass string `json:"pass,omitempty"`
}

type Logging struct {
	Dir string `json:"dir"`
	Level int `json:"level"`
}

type AppConfig struct {
	Env    string   `json:"env"`
	Debug  string   `json:"debug"`
	Domain string   `json:"domian"`
	Port   int      `json:"port"`
	Mysql  DB `json:"mysql"`
	Mongo  DB `json:"mongo"`
	Redis  struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		DB   int    `json:"db"`
	} `json:"redis"`
	Log Logging `json:"log"`
}

var Config AppConfig

func Load(env string) (*AppConfig, error) {
	configFile, err := os.Open(env + ".config.json")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)
	if err != nil {
		return nil, err
	}

	return &Config, nil
}
