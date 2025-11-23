package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	RedisAddress  string `envconfig:"REDIS_ADDR" default:"localhost:6379"`
	RedisPassword string `envconfig:"REDIS_PWD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`
	Port          string `envconfig:"PORT" default:"8080"`
}

func EnvParse() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
