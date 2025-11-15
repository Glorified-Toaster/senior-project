package main

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/mongodb"
	"github.com/Glorified-Toaster/senior-project/internal/config/logger"
	"github.com/Glorified-Toaster/senior-project/internal/controllers"
	"github.com/Glorified-Toaster/senior-project/internal/repository"
	"github.com/Glorified-Toaster/senior-project/internal/server"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/go-playground/validator"
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

	// init cache
	cache, err := cache.InitCache("myapp")
	if err != nil {
		utils.LogErrorWithLevel("error", utils.DragonflyFailedToInit.Type, utils.DragonflyFailedToInit.Code, utils.DragonflyFailedToInit.Msg, err)
	}
	// init the user repo
	repo := repository.NewUserRepo(context.Background(), mongodb.Database, cache)
	// init validator
	validate := validator.New()
	// pass cache, repo , validator to controllers
	controllers.NewControllers(validate, *repo, *cache)

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
