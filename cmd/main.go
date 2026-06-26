package main

import (
	"SpotSync/internal/config"
	"SpotSync/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db := config.ConnectDatabase(cfg)

	server.StartServer( db, cfg)
}
