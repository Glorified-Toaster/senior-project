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
