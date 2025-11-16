package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RedisAddress  string `envconfig:"" default:"localhost:6379"`
	RedisPassword string `envconfig:"" default:""`
	RedisDB       int    `envconfig:"" default:"0"`
	Port          string `envconfig:"" default:"8080"`
}

func EnvParse() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
