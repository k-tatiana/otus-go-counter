package main

import (
	"go-server-counters/config"
	"go-server-counters/server"
)

func main() {
	cfg, err := config.EnvParse()
	if err != nil {
		panic(err)
	}

	server.Run(cfg)
}
