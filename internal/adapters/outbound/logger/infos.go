package logger

type Info struct {
	Type string
	Msg  string
}

// Info types
const (
	DatabaseInfo       string = "DATABASE_INFO"
	CacheInfo          string = "CACHE_INFO"
	InternalServerInfo string = "INTERNAL_INFO"
)

// info
var (
	// mongodb info
	PostgresIsConnected = Info{
		DatabaseInfo,
		"Connected to PostgreSQL successfully...",
	}

	PostgresIsDisconnected = Info{
		DatabaseInfo,
		"Disconnected from mangodb successfully...",
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
