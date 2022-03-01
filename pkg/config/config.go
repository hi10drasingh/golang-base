package config

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	DB   string `json:"db"`
	User string `json:"user,omitempty"`
	Pass string `json:"pass,omitempty"`
}

type LogConfig struct {
	Dir string `json:"dir"`
}

type AppConfig struct {
	Env    string   `json:"env"`
	Debug  string   `json:"debug"`
	Domain string   `json:"domian"`
	Port   int      `json:"port"`
	Mysql  DBConfig `json:"mysql"`
	Mongo  DBConfig `json:"mongo"`
	Redis  struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		DB   int    `json:"db"`
	} `json:"redis"`
	Log LogConfig `json:"log"`
}

var Config AppConfig

func Load(env string) (*AppConfig, error) {
	configFile, err := os.Open(env + ".config.json")
	if err != nil {
		return &Config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)
	if err != nil {
		return &Config, err
	}

	return &Config, nil
}
