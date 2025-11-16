package main

import (
	"go-server-counters/config"
)

func main() {
	cfg, err := config.EnvParse()
	if err != nil {
		panic(err)
	}

	runServer(cfg)
}
