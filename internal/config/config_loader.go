// Package config provides functionality to initialize and
// load configuration from YAML file using viper package.
package config

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	GlobalConfig *Config
	once         sync.Once
)

type Config struct {
	HTTPServer  *HTTPServerConf  `yaml:"http_server" mapstructure:"http_server"`
	MongoDB     *MongoDBConf     `yaml:"mongodb" mapstructure:"mongodb"`
	DragonflyDB *DragonflyDBConf `yaml:"dragonflydb" mapstructure:"dragonflydb"`
}

type HTTPServerConf struct {
	Addr     string `yaml:"address" mapstructure:"address"`
	Port     string `yaml:"port" mapstructure:"port"`
	CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `yaml:"key_file" mapstructure:"key_file"`
	CertDir  string `yaml:"tls_cert_dir" mapstructure:"tls_cert_dir"`
}

type MongoDBConf struct {
	Database string `yaml:"database" mapstructure:"database"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
}

type DragonflyDBConf struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

// Init : to initialize the configuration loading process.
func Init(path, file string) {
	var err error

	// to ensure that the config is loaded only once
	once.Do(func() {
		// flag definitions
		configPath := flag.String("config-path", path, "path to configuration path")
		configFile := flag.String("config-file", file, "name of configuration file (without extension)")
		isDebugMode := flag.Bool("debug", false, "enable gin debug mode")

		if !flag.Parsed() {
			flag.Parse()
		}

		if *isDebugMode {
			gin.SetMode(gin.DebugMode)
			log.Println("Server is running in debug mode.")
		} else {
			gin.SetMode(gin.ReleaseMode)
			log.Println("Server is running in release mode.")
		}

		GlobalConfig, err = LoadConfig(*configPath, *configFile)
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		log.Println("Configuration loaded successfully")
	})
}

// LoadConfig : to load the YAML configuration file (using viper package).
func LoadConfig(configPath, configFile string) (*Config, error) {
	var config *Config

	viperInst := viper.New() // init viper instance
	viperInst.AddConfigPath(configPath)
	viperInst.SetConfigName(configFile)
	viperInst.SetConfigType("yaml")

	setDefaultConfig(viperInst)

	if err := viperInst.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	if err := viperInst.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshaling config error : %w", err)
	}

	return config, nil
}

// GetConfig : to get the global configuration instance.
func GetConfig() (*Config, error) {
	if GlobalConfig == nil {
		return nil, fmt.Errorf("configuration not initialized you must call config.Init() first")
	}
	return GlobalConfig, nil
}

// setDefaultConfig : to set important default values.
func setDefaultConfig(viperInst *viper.Viper) {
	// setting important default values
	viperInst.SetDefault("http_server.address", "localhost")
	viperInst.SetDefault("http_server.port", "8080")
	viperInst.SetDefault("http_server.tls_cert_dir", "certs")

	// MongoDB default values
	viperInst.SetDefault("mongodb.host", "localhost")
	viperInst.SetDefault("mongodb.port", "27017")
	viperInst.SetDefault("mongodb.database", "senior_project")

	// DragonflyDB default values
	viperInst.SetDefault("dragonflydb.host", "localhost")
	viperInst.SetDefault("dragonflydb.port", "6379")
	viperInst.SetDefault("dragonflydb.db", 0)
}
