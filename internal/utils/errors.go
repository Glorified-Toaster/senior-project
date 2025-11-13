// Package utils provides most used variables and functionalities.
package utils

import (
	"go.uber.org/zap"
)

// error types
const (
	DatabaseError       string = "DATABASE_ERROR" // for mongo errors
	CacheError          string = "CACHE_ERROR"    // for Dragonfly errors
	InternalServerError string = "INTERNAL_ERROR" // for internal errors
)

type Error struct {
	Type string
	Code string
	Msg  string
}

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

	MongoNotInitialized = Error{
		DatabaseError,
		"MONGODB_NOT_INITIALIZED_ERROR",
		"mongodb client is not initialized",
	}

	MongoFailedToGetCollection = Error{
		DatabaseError,
		"MONGODB_FAILED_TO_GET_COLLECTION_ERROR",
		"failed to get collection",
	}

	// dragonfly errors
	DragonflyFailedToInit = Error{
		CacheError,
		"DRAGONFLYDB_CONNECTION_ERROR",
		"failed to init dragonflydb",
	}

	DragonflyFailedToLoadOptions = Error{
		CacheError,
		"DRAGONFLYDB_CONFIG_LOAD_ERROR",
		"unable to load dragonfly config options",
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

	FailedToGenerateTLSCert = Error{
		InternalServerError,
		"FAILED_TO_GENERATE_TLS_ERROR",
		"failed to generate self-signed TLS certificate",
	}

	FailedToStartWithTLS = Error{
		InternalServerError,
		"FAILED_TO_START_WITH_TLS_ERROR",
		"unable to use TLS with HTTPS/2, using HTTP/1.1 instead",
	}

	FailedToStartServer = Error{
		InternalServerError,
		"FAILED_TO_START_SERVER_ERROR",
		"unable to start http server",
	}

	FailedToShutdownServer = Error{
		InternalServerError,
		"FAILED_TO_SHUTDOWN_SERVER_ERROR",
		"failed to shutdown the http server",
	}
)

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

	case "panic":
		zlogger.Panic(errorType, allFields...)

	default:
		zlogger.Error(errorType, allFields...)

	}
}
