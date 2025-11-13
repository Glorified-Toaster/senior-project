package utils

import "go.uber.org/zap"

// Info types
const (
	DatabaseInfo       string = "DATABASE_INFO"
	CacheInfo          string = "CACHE_INFO"
	InternalServerInfo string = "INTERNAL_INFO"
)

type Info struct {
	Type string
	Msg  string
}

// info
var (

	// mongodb info
	MongoIsConnected = Info{
		DatabaseInfo,
		"Connected to mongodb successfully...",
	}

	MongoIsDisconnected = Info{
		DatabaseInfo,
		"Disconnected from mangodb successfully...",
	}

	// dragonfly info
	DragonflyIsConnected = Info{
		CacheInfo,
		"Connected to dragonflydb successfully...",
	}

	// internal info
	ServerStartOK = Info{
		InternalServerInfo,
		"Starting the go server",
	}

	GenerateTLSCertOK = Info{
		InternalServerInfo,
		"Self-signed TLS generated successfully...",
	}

	UseExistedTLSCert = Info{
		InternalServerInfo,
		"Using an existing TLS certificate...",
	}

	ServerShutdownSignalOK = Info{
		InternalServerInfo,
		"Shutdown signal recived",
	}

	ServerShutdown = Info{
		InternalServerInfo,
		"Shutting down the server...",
	}

	ServerShutdownOK = Info{
		InternalServerInfo,
		"HTTP server shutdown gracefully...",
	}
)

// LogInfo : log Info (very useful comment i guess...)
func LogInfo(infoType, msg string, fields ...zap.Field) {
	mandatoryFields := []zap.Field{
		zap.String("info_msg", msg),
	}

	allFields := append(mandatoryFields, fields...)
	zlogger.Info(infoType, allFields...)
}
