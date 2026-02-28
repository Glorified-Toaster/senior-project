package main

import (
	"context"
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
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/adapters/outbound/logger"
	"uot-exam/internal/adapters/outbound/repository"
	"uot-exam/internal/application"
	"uot-exam/internal/ports"

	"github.com/bxcodec/faker/v4"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

	DBConfig := &database.DBConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DatabaseName,
		SSLMode:         cfg.Database.SSLMode,
		MaxConns:        cfg.Database.MaxConns,
		MinConns:        cfg.Database.MinConns,
		MaxConnLifetime: cfg.Database.MaxConnLifetime,
		MaxConnIdleTime: cfg.Database.MaxConnIdleTime,
	}

	pool, err := database.NowConnection(*DBConfig)
	if err != nil {
		zlog.LogErrorWithLevel("fatal", logger.DatabaseError, logger.MongoFailedToConnect.Code, "Database failed to connect", err)
		return
	}
	zlog.LogInfo(logger.MongoIsConnected.Type, logger.MongoIsConnected.Msg)

	query := sqlc.New(pool)
	userRepo := repository.NewUserRepository(query)
	txManager := database.NewPostgresTxManager(pool.Pool)
	app := application.NewApplication(userRepo, txManager, pool, zlog)

	//populateDB(app)
	//mockExam(app)

	pool.Stats()
	// init validator
	validate := validator.New()
	// init jwt
	jwt := helpers.NewJWT(cfg)
	// init auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwt, zlog)
	// pass cache, repo, validator, jwt to controllers
	userCtrl := handler.NewUserHandler(app, validate, jwt, cfg, zlog)

	// initialize the server
	srv := server.NewServer(userCtrl, authMiddleware, zlog, cfg)

	zlog.LogInfo(logger.ServerStartOK.Type, logger.ServerStartOK.Msg, zap.String("server_address", net.JoinHostPort(cfg.HTTPServer.Addr, cfg.HTTPServer.Port)))

	// start the server over TLS
	srv.StartOverTLS(cfg)
}

func mockExam(app *application.Application) {
	user, err := app.GetUserByID(context.Background(), uuid.MustParse("69744b6d-0101-411e-b617-bbea7ade4b74"))
	if err != nil {
		log.Println("error getting user by id", err)
	}
	log.Println("user", user)

	user, err = app.GetUserByUsername(context.Background(), "Lind9348")
	if err != nil {
		log.Println("error getting user by username", err)
	}
	log.Println("user", user)

}

func populateDB(app *application.Application) {
	for i := 0; i < 100; i++ {
		user, err := app.CreateUser(context.Background(), ports.CreateUserParams{
			Username: faker.Username(),
			FullName: faker.Name(),
			Password: faker.Password(),
			Role:     "STUDENT",
			IsActive: true,
		})
		if err != nil {
			log.Println("error creating user", err)
		}
		log.Println("user", user)
	}
}
