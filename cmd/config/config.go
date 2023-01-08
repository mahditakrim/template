package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		LogPath   string    `yaml:"log_path"`
		Transport transport `yaml:"transport"`
		DB        db        `yaml:"db"`
	}

	transport struct {
		HttpAddr string `yaml:"http_addr"`
		RpcAddr  string `yaml:"rpc_addr"`
	}

	db struct {
		Postgres struct {
			Addr     string `yaml:"addr"`
			Db       string `yaml:"db"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"postgres"`
	}
)

var globalConfig Config

func Init() error {

	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, &globalConfig)
}

func Get() Config {

	return globalConfig
}
