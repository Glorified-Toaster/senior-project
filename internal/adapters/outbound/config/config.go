package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Configuration struct {
	Config *Config
}

type Config struct {
	HTTPServer  *HTTPServerConf   `yaml:"http_server" mapstructure:"http_server"`
	Database    *DatabaseConf     `yaml:"database" mapstructure:"database"`
	DragonflyDB *DragonflyDBConf  `yaml:"dragonflydb" mapstructure:"dragonflydb"`
	ZapLogger   *ZapLoggerConf    `yaml:"zap_logger" mapstructure:"zap_logger"`
	Lumberjack  *LumberjackConf   `yaml:"lumberjack" mapstructure:"lumberjack"`
	JWTAuth     *JWTAuthConf      `yaml:"jwt_auth" mapstructure:"jwt_auth"`
	CSRF        *CSRFConfig       `yaml:"csrf" mapstructure:"csrf"`
	GinSession  *GinSessionConfig `yaml:"gin_session" mapstructure:"gin_session"`
	GinLogger   *GinLoggerConf    `yaml:"gin_logger" mapstructure:"gin_logger"`
}

func (cfg *Configuration) Init(path, file string) error {
	configPath, configFile := setFlags(path, file)

	config, err := cfg.loadConfig(configPath, configFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	cfg.Config = config

	return nil
}

func setFlags(path, file string) (*string, *string) {
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

	return configPath, configFile
}

func (cfg *Configuration) loadConfig(path, file *string) (*Config, error) {
	viperInst := viper.New()
	viperInst.AddConfigPath(*path)
	viperInst.SetConfigName(*file)
	viperInst.SetConfigType("yaml")

	setDefaultConfig(viperInst)

	if err := viperInst.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	config := &Config{}

	if err := viperInst.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unmarshaling config error : %w", err)
	}

	return config, nil
}

func (cfg *Configuration) GetConfig() (*Config, error) {
	if cfg == nil || cfg.Config == nil {
		return nil, fmt.Errorf("configuration not initialized, call Init() first")
	}
	return cfg.Config, nil
}

func (cfg *Config) Validate() error {
	if cfg.HTTPServer == nil {
		return fmt.Errorf("http_server configuration is required")
	}
	if cfg.Database == nil {
		return fmt.Errorf("database configuration is required")
	}
	if cfg.ZapLogger == nil {
		return fmt.Errorf("logger configuration is required")
	}
	if cfg.Lumberjack == nil {
		return fmt.Errorf("lumberjack configuration is required")
	}
	if cfg.DragonflyDB == nil {
		return fmt.Errorf("dragonflydb configuration is required")
	}
	if cfg.JWTAuth == nil {
		return fmt.Errorf("jwt_auth configuration is required")
	}
	if cfg.GinSession == nil {
		return fmt.Errorf("gin_session configuration is required")
	}
	if cfg.CSRF == nil {
		return fmt.Errorf("csrf configuration is required")
	}
	if cfg.GinLogger == nil {
		return fmt.Errorf("gin_logger configuration is required")
	}

	return nil
}
