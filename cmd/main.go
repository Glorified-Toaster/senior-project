package main

import (
	"fmt"

	"github.com/Glorified-Toaster/senior-project/internal/config"
)

func main() {
	// loading the env variables
	config.Init("/home/potato/Dev/go/senior-project/internal/config", "config")
	cfg := config.GlobalConfig
	fmt.Println(cfg.HTTPServer.Port)
	fmt.Println(cfg.HTTPServer.Addr)
}
