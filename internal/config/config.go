// Package config provides loading and managing configuration variables.
// via YAML file and env variables. through hot-reloding and values override.
package config

import (
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	GlobalConfig *Config
	configMutex  sync.RWMutex
)

type Config struct {
	HTTPServer *HTTPServerConf `yaml:"http_server" mapstructure:"http_server"`
}

type HTTPServerConf struct {
	Addr string `yaml:"address" mapstructure:"address"`
	Port string `yaml:"port" mapstructure:"port"`
}

// Init : init with cmd flags or using the passed params values
func Init(path, file string) {
	configPath := flag.String("config-path", path, "path to configuration path")
	configFile := flag.String("config-file", file, "name of configuration file (without extension)")
	flag.Parse()

	conf, err := LoadConfig(*configPath, *configFile)
	if err != nil {
		panic("Load config fail : " + err.Error())
	}
	GlobalConfig = conf
}

// LoadConfig : to load the YAML configuration file (using viper package).
func LoadConfig(configPath, configFile string) (*Config, error) {
	var config *Config

	viperInst := viper.New() // init viper instance
	viperInst.AddConfigPath(configPath)
	viperInst.SetConfigName(configFile)

	viperInst.SetConfigType("yaml")

	// reading env variables from file and handle the error.
	if err := viperInst.ReadInConfig(); err != nil {
		return nil, err
	}

	/*
		enable env variables to override config:
		example : if port:8080 and an env variable exported via :

				export WEB_APP_PORT=9090

		then the port would be port:9090
	*/
	viperInst.SetEnvPrefix("WEB_APP")
	viperInst.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperInst.AutomaticEnv()

	if err := viperInst.Unmarshal(&config); err != nil {
		return nil, err
	}

	envOverrides(config)

	configMutex.Lock()
	GlobalConfig = config
	configMutex.Unlock()

	// file change monitoring
	viperInst.WatchConfig()
	viperInst.OnConfigChange(func(in fsnotify.Event) {
		var newConfig Config

		if err := viperInst.Unmarshal(&newConfig); err == nil {
			envOverrides(&newConfig)

			configMutex.Lock()
			GlobalConfig = &newConfig
			configMutex.Unlock()

		}
	})
	return config, nil
}

// envOverrides : apply env variables overrides.
func envOverrides(config *Config) {
	httpEnvOverrides(config)
}

// httpEnvOverrides : apply env overrides for http related variables.
func httpEnvOverrides(config *Config) {
	if addr := os.Getenv("WEB_APP_HTTP_SERVER_ADDR"); addr != "" {
		config.HTTPServer.Addr = addr
	}

	if port := os.Getenv("WEB_APP_HTTP_SERVER_PORT"); port != "" {
		config.HTTPServer.Port = port
	}
}
