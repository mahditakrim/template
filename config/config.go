package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		LogPath string `yaml:"log_path"`
		Web     web    `yaml:"web"`
		DB      db     `yaml:"db"`
	}

	web struct {
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

func Init() (*Config, error) {

	yamlFile, err := os.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(yamlFile, conf)

	return conf, err
}
