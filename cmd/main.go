package main

import (
	"gotickets/internal/config"
	"gotickets/internal/server"
)

// @title           GoTickets API
// @version         1.0
// @description     A professional backend service for a ticket booking platform.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@gotickets.local

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:5000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// load environment variables
	cfg := config.LoadEnv()
	// connect to the database
	db := config.ConnectDatabase(cfg)
	// start the server
	server.Start(db, cfg)

}
