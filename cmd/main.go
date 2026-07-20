package main

import (
	"gotickets/internal/config"
	"gotickets/internal/server"
)

func main() {
	// load environment variables
	cfg := config.LoadEnv()
	// connect to the database
	db := config.ConnectDatabase(cfg)
	// start the server
	server.Start(db, cfg)

}
