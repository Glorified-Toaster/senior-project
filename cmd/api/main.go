package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"

	"uot-exam/internal/adapters/inbound/http/handler"
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/inbound/http/middleware"
	"uot-exam/internal/adapters/inbound/http/server"
	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/adapters/outbound/logger"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

func main() {
	// Initialize Configuration
	configInstance := &config.Configuration{}
	if err := configInstance.Init("/home/potato/Dev/final-project/config", "config"); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	cfg, err := configInstance.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get configuration: %v", err)
	}

	// Initialize Zap
	zapLogger, err := logger.InitZapLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		err := zapLogger.Sync()
		if err != nil && !errors.Is(err, syscall.EINVAL) {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}()

	// Initialize Logger
	zlog := logger.New(zapLogger)

	// init validator
	validate := validator.New()
	// init jwt
	jwt := helpers.NewJWT(cfg)
	// init auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwt, zlog)
	// pass cache, repo, validator, jwt to controllers
	ctrl := handler.NewHandler(validate, jwt, cfg, zlog)

	// initialize the server
	srv := server.NewServer(ctrl, authMiddleware, zlog)

	zlog.LogInfo(logger.ServerStartOK.Type, logger.ServerStartOK.Msg, zap.String("server_address", net.JoinHostPort(cfg.HTTPServer.Addr, cfg.HTTPServer.Port)))

	// start the server over TLS
	srv.StartOverTLS(cfg)
}
