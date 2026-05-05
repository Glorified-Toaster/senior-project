package logger

type Error struct {
	Type string
	Code string
	Msg  string
}

// error types
const (
	DatabaseError       string = "DATABASE_ERROR" // for postgres errors
	CacheError          string = "CACHE_ERROR"    // for Dragonfly errors
	InternalServerError string = "INTERNAL_ERROR" // for internal errors
)

// errors
var (

	// postgres errors
	PosFailedToConnect = Error{
		DatabaseError,
		"POSTGRES_CONNECTION_ERROR",
		"failed to connect to postgres",
	}
	PosFailedToDisconnect = Error{
		DatabaseError,
		"POSTGRES_DISCONNECTION_ERROR",
		"failed to disconnect from postgres",
	}

	PosNotInitialized = Error{
		DatabaseError,
		"POSTGRES_NOT_INITIALIZED_ERROR",
		"postgres client is not initialized",
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
