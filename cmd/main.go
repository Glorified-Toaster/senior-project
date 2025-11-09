package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/server"
)

func main() {
	// loading the env variables
	config.Init(getConfigPath(), "config")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
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
