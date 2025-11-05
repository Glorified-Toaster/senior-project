package main

import (
	"log"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/server"
)

func main() {
	// loading the env variables
	config.Init("/home/potato/Dev/go/senior-project/internal/config", "config")

	srv := server.NewServer()

	log.Println("Starting server on " + config.GlobalConfig.HTTPServer.Addr + ":" + config.GlobalConfig.HTTPServer.Port)
	srv.StartOverTLS("certs/cert.pem", "certs/key.pem")
}
