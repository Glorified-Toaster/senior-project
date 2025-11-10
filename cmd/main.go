package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/mongodb"
	"github.com/Glorified-Toaster/senior-project/internal/server"
)

// --- AI GEN ---//
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// --- AI GEN ---//

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

	cache, err := cache.InitCache("myapp")
	if err != nil {
		log.Fatal(err)
	}

	// --- AI GEN ---//
	if err := cache.HealthCheck(); err != nil {
		log.Printf("Cache health check failed: %v", err)
	} else {
		log.Println("Cache health check is : OK")
	}

	// Example usage
	user := User{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Set value
	err = cache.Set("user:123", user, 10*time.Minute)
	if err != nil {
		log.Printf("Error setting cache: %v", err)
	}

	// Get value
	var retrievedUser User
	err = cache.Get("user:123", &retrievedUser)
	if err != nil {
		log.Printf("Error getting cache: %v", err)
	} else {
		log.Printf("Retrieved user: %+v", retrievedUser)
	}
	// --- AI GEN ---//

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
