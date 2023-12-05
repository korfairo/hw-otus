package main

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Memory      bool
	PostgresDSN string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig(configFilePath string) (Config, error) {
	config := Config{}
	_, err := toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return Config{}, errors.Wrap(err, "failed to decode toml")
	}
	return config, nil
}
