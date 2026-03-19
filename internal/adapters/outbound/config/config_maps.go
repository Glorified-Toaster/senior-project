package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	// HTTP Server
	DefaultHTTPAdress = "localhost"
	DefaultHTTPPort   = "8080"
	DefaultTLSDir     = "certs"

	// database
	DefaultDBHost = "localhost"
	DefaultDBPort = "5432"
	DefaultDBName = "UOT-MCQ-Exam"
)

type HTTPServerConf struct {
	Addr     string `yaml:"address" mapstructure:"address"`
	Port     string `yaml:"port" mapstructure:"port"`
	CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `yaml:"key_file" mapstructure:"key_file"`
	CertDir  string `yaml:"tls_cert_dir" mapstructure:"tls_cert_dir"`
}

type DatabaseConf struct {
	DatabaseName    string        `yaml:"database_name" mapstructure:"database_name"`
	Username        string        `yaml:"username" mapstructure:"username"`
	Password        string        `yaml:"password" mapstructure:"password"`
	Host            string        `yaml:"host" mapstructure:"host"`
	Port            int           `yaml:"port" mapstructure:"port"`
	SSLMode         string        `yaml:"ssl_mode" mapstructure:"ssl_mode"`
	MaxConns        int           `yaml:"max_connections" mapstructure:"max_connections"`
	MinConns        int           `yaml:"min_connections" mapstructure:"min_connections"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" mapstructure:"max_conn_idle_time"`
}

type DragonflyDBConf struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type ZapLoggerConf struct {
	Development bool   `yaml:"development" mapstructure:"development"`
	Level       string `yaml:"level" mapstructure:"level"`
	Encoding    string `yaml:"encoding" mapstructure:"encoding"`
	DirPath     string `yaml:"log_dir" mapstructure:"log_dir"`
	FileName    string `yaml:"log_file" mapstructure:"log_file"`
}

type LumberjackConf struct {
	MaxSize    int  `yaml:"max_size" mapstructure:"max_size"`
	MaxAge     int  `yaml:"max_age" mapstructure:"max_age"`
	MaxBackups int  `yaml:"max_backups" mapstructure:"max_backups"`
	Compress   bool `yaml:"compress" mapstructure:"compress"`
}

type JWTAuthConf struct {
	Secret string `yaml:"secret" mapstructure:"secret"`
}

type CSRFConfig struct {
	Secret string `yaml:"secret" mapstructure:"secret"`
}

type GinSessionConfig struct {
	Secret string `yaml:"secret" mapstructure:"secret"`
}

type GinLoggerConf struct {
	Filename string `yaml:"filename" mapstructure:"filename"`
}

// setDefaultConfig : to set important default values.
func setDefaultConfig(viperInst *viper.Viper) {
	// setting important default values
	viperInst.SetDefault("http_server.address", DefaultHTTPAdress)
	viperInst.SetDefault("http_server.port", DefaultHTTPPort)
	viperInst.SetDefault("http_server.tls_cert_dir", DefaultTLSDir)

	// Database default values
	viperInst.SetDefault("database.host", DefaultDBHost)
	viperInst.SetDefault("database.port", DefaultDBPort)
	viperInst.SetDefault("database.database_name", DefaultDBName)

	// DragonflyDB default values
	viperInst.SetDefault("dragonflydb.host", "localhost")
	viperInst.SetDefault("dragonflydb.port", "6379")
	viperInst.SetDefault("dragonflydb.db", 0)

	// Zap default values
	viperInst.SetDefault("zap_logger.log_dir", "./logs")
	viperInst.SetDefault("zap_logger.development", false)
	viperInst.SetDefault("zap_logger.level", "debug")
	viperInst.SetDefault("zap_logger.encoding", "json")
	viperInst.SetDefault("zap_logger.log_file", "app.log")

	// Gin default values
	viperInst.SetDefault("gin_logger.filename", "logs/gin.log")
}
