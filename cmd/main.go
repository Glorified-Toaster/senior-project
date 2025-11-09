package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/mongodb"
	"github.com/Glorified-Toaster/senior-project/internal/server"
)

func main() {
	// loading the env variables
	config.Init(getConfigPath(), "config")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	// connect to MongoDB
	if cfg.MongoDB != nil {
		uri := mongodb.MakeURI(
			cfg.MongoDB.Host,
			cfg.MongoDB.Port,
			cfg.MongoDB.Username,
			cfg.MongoDB.Password,
			cfg.MongoDB.Database,
		)

		if err := mongodb.MongoConnect(uri, cfg.MongoDB.Database); err != nil {
			log.Fatalf("failed to connect to MongoDB: %v", err)
		}
		defer func() {
			if err := mongodb.MongoDisconnect(); err != nil {
				log.Printf("error disconnecting from MongoDB: %v", err)
			}
		}()
	}

	// initialize the server
	srv := server.NewServer()

	log.Println("Starting server on " + net.JoinHostPort(cfg.HTTPServer.Addr, cfg.HTTPServer.Port) + "...")
	// start the server over TLS
	srv.StartOverTLS(cfg.HTTPServer.CertFile, cfg.HTTPServer.KeyFile)
}

// getConfigPath : to get the path to the config file.
func getConfigPath() string {
	return filepath.Join(getProgramPath(), "internal", "config")
}

// getProgramPath : to get the path where the program is located.
func getProgramPath() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable path: %v", err)
	}
	return filepath.Dir(exe)
}
