package main

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/mongodb"
	"github.com/Glorified-Toaster/senior-project/internal/config/logger"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/repository"
	"github.com/Glorified-Toaster/senior-project/internal/server"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func main() {
	// loading the YAML config variables
	config.Init(getConfigPath(), "config")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("[%s]-%s : %v", utils.ConfigFailedToLoad.Code, utils.ConfigFailedToLoad.Msg, err)
	}

	// init zap logger
	err = logger.InitLogger(*cfg)
	if err != nil {
		log.Fatalf("[%s]-%s : %v", utils.LoggerFailedToInit.Code, utils.LoggerFailedToInit.Msg, err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("[%s]-%s : %v", utils.LoggerFailedToSync.Code, utils.LoggerFailedToSync.Msg, err)
		}
	}()

	// get logger for the utils package
	utils.InitUtils()

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
			utils.LogErrorWithLevel("fatal",
				utils.MongoFailedToConnect.Type,
				utils.MongoFailedToConnect.Code,
				utils.MongoFailedToConnect.Msg,
				err,
			)
		}
		defer func() {
			if err := mongodb.MongoDisconnect(); err != nil {
				utils.LogErrorWithLevel("warn",
					utils.MongoFailedToDisconnect.Type,
					utils.MongoFailedToDisconnect.Code,
					utils.MongoFailedToDisconnect.Msg,
					err,
				)
			}
		}()
	}

	cache, err := cache.InitCache("myapp")
	if err != nil {
		utils.LogErrorWithLevel("error", utils.DragonflyFailedToInit.Type, utils.DragonflyFailedToInit.Code, utils.DragonflyFailedToInit.Msg, err)
	}

	// ---- Test ----

	repo := repository.NewUserRepo(mongodb.Database, cache)

	user := createExampleUser()
	repo.CreateUser(context.Background(), user)
	repo.GetUserByID(context.Background(), "691745daf91a3e3357b035e2")

	// ---- Test ----

	// initialize the server
	srv := server.NewServer()

	utils.LogInfo(utils.ServerStartOK.Type, utils.ServerStartOK.Msg, zap.String("server_address", net.JoinHostPort(cfg.HTTPServer.Addr, cfg.HTTPServer.Port)))

	// start the server over TLS
	srv.StartOverTLS(cfg)
}

// getConfigPath : to get the path to the config file.
func getConfigPath() string {
	return filepath.Join(getProgramPath(), "internal", "config")
}

// getProgramPath : to get the path where the program is located.
func getProgramPath() string {
	exe, err := os.Executable()
	if err != nil {
		utils.LogErrorWithLevel("fatal", utils.FailedToGetProgramPath.Type, utils.FailedToGetProgramPath.Code, utils.FailedToGetProgramPath.Msg, err)
	}
	return filepath.Dir(exe)
}

// ---- Test ----

func createExampleUser() *models.User {
	id := primitive.NewObjectID()

	firstName := "John"
	lastName := "Doe"
	password := "securePassword123"
	email := "john.doe@example.com"
	phone := "+1-555-0123"
	token := "auth_token_abc123"
	role := "USER"
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

	return &models.User{
		ID:          id,
		FirstName:   &firstName,
		LastName:    &lastName,
		Password:    &password,
		Email:       &email,
		Phone:       &phone,
		Token:       &token,
		Role:        &role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      id.Hex(),
		AccessToken: accessToken,
	}
}

// ---- Test ----
