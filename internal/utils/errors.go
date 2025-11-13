// Package utils provides most used variables and functionalities.
package utils

import (
	"sync"

	"github.com/Glorified-Toaster/senior-project/internal/config/logger"
	"go.uber.org/zap"
)

var (
	zlogger *zap.Logger
	once    sync.Once
)

// error types
const (
	DatabaseError       string = "DATABASE_ERROR" // for mongo errors
	CacheError          string = "CACHE_ERROR"    // for Dragonfly errors
	InternalServerError string = "INTERNAL_ERROR" // for internal errors
)

// Info types
const (
	DatabaseInfo       string = "DATABASE_INFO"
	CacheInfo          string = "CACHE_INFO"
	InternalServerInfo string = "INTERNAL_INFO"
)

type Error struct {
	Type string
	Code string
	Msg  string
}

type Info struct {
	Type string
	Msg  string
}

// info
var (
	ServerStart = Info{
		InternalServerInfo,
		"Starting the go server",
	}
)

// errors
var (

	// mongo errors
	MongoFailedToConnect = Error{
		DatabaseError,
		"MONGODB_CONNECTION_ERROR",
		"failed to connect to mongodb",
	}
	MongoFailedToDisconnect = Error{
		DatabaseError,
		"MONGODB_DISCONNECTION_ERROR",
		"failed to disconnect from mongodb",
	}

	// dragonfly errors
	DragonflyFailedToInit = Error{
		CacheError,
		"DRAGONFLYDB_CONNECTION_ERROR",
		"failed to init dragonflydb",
	}

	// internal errors
	ConfigFailedToLoad = Error{
		InternalServerError,
		"CONFIG_LOAD_ERROR",
		"failed to get config",
	}

	LoggerFailedToInit = Error{
		InternalServerError,
		"LOGGER_INIT_ERROR",
		"failed to init zap logger",
	}

	LoggerFailedToSync = Error{
		InternalServerError,
		"LOGGER_SYNC_ERROR",
		"failed to sync zap logger",
	}

	FailedToGetProgramPath = Error{
		InternalServerError,
		"PROGRAM_PATH_LOAD_ERROR",
		"failed to get executable path",
	}
)

// InitUtils : get logger (once)
func InitUtils() {
	once.Do(func() {
		zlogger = logger.GetLogger()
	})
}

// LogErrorWithLevel : log error and select the level of that error
func LogErrorWithLevel(level, errorType, errorCode, msg string, err error, fields ...zap.Field) {
	// must be passed fields
	mandatoryFields := []zap.Field{
		zap.String("error_code", errorCode),
		zap.String("error_msg", msg),
		zap.Error(err),
	}

	// append all of the additional fields
	allFields := append(mandatoryFields, fields...)

	switch level {
	case "fatal":
		zlogger.Fatal(errorType, allFields...)
	case "error":
		zlogger.Error(errorType, allFields...)
	case "warn":
		zlogger.Warn(errorType, allFields...)
	}
}

// LogInfo : log Info (very useful comment i guess...)
func LogInfo(infoType, msg string, fields ...zap.Field) {
	mandatoryFields := []zap.Field{
		zap.String("info_msg", msg),
	}

	allFields := append(mandatoryFields, fields...)
	zlogger.Info(infoType, allFields...)
}
