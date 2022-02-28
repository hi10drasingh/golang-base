package config

import (
	"encoding/json"
	"flag"
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

var env string = "local"

var Config AppConfig

func init() {
	// Path to config file can be passed in.
	flag.StringVar(&env, "env", env, "Environment")
	flag.Parse()

	configFile, err := os.Open(env + ".config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)
	if err != nil {
		panic(err)
	}

}
