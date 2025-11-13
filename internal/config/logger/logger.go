// Package logger for logs with zap and lumberjack for log rotation
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Glorified-Toaster/senior-project/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ZapLogger *zap.Logger
	once      sync.Once
)

// InitLogger : initializes the Zap logger once
func InitLogger(cfg config.Config) error {
	var initErr error
	once.Do(func() {
		ZapLogger, initErr = LogWithZap(cfg)
	})
	return initErr
}

// GetLogger : return a zap log instance
func GetLogger() *zap.Logger {
	// fallback measure
	if ZapLogger == nil {
		logger, _ := zap.NewProduction()
		return logger
	}
	return ZapLogger
}

// Sync : calls the underlying Core's Sync method, flushing any buffered log entries.
func Sync() error {
	if ZapLogger != nil {
		err := ZapLogger.Sync()
		// ignore stdout error because stdout doesn't support sync anyway.
		if err != nil && strings.Contains(err.Error(), "invalid argument") {
			return nil
		}
		return err
	}
	return nil
}

// LogWithZap : create a zap instance and lumberjack for file rotation
func LogWithZap(cfg config.Config) (*zap.Logger, error) {
	logDir := cfg.ZapLogger.DirPath
	// check if dir exsists if not make one
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		log.Printf("Creating log directory : %s\n", logDir)
		if err := os.MkdirAll(logDir, 0o755); // rwxr-xr-x permission
		err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to access log directory: %v", err)
	}

	logFile := filepath.Join(logDir, cfg.ZapLogger.FileName)

	// create lumberjack logger config struct
	lumberjack := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.Lumberjack.MaxSize,
		MaxBackups: cfg.Lumberjack.MaxBackups,
		MaxAge:     cfg.Lumberjack.MaxAge,
		Compress:   cfg.Lumberjack.Compress,
	}

	// determine what level will zap show [Debug(default) , Info , Warn , Error]
	level, err := zapcore.ParseLevel(cfg.ZapLogger.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid zap level : %v", err)
	}

	// how logs will be encoded (default is production (json))
	var zapEncodeCfg zapcore.EncoderConfig

	if cfg.ZapLogger.Development {
		zapEncodeCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		zapEncodeCfg = zap.NewProductionEncoderConfig()
	}

	// ISO8601 : YYYY-MM-DDTHH:MM:SS.milliseconds(Timezone offset from UTC)
	zapEncodeCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// make a console encoder and json encoder
	consoleEncoder := zapcore.NewConsoleEncoder(zapEncodeCfg)
	fileEncoder := zapcore.NewJSONEncoder(zapEncodeCfg)

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjack), level)

	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger.Info("Logger initialized successfully",
		zap.String("level", cfg.ZapLogger.Level),
		zap.String("log_file", logFile),
		zap.Bool("development", cfg.ZapLogger.Development),
	)

	return logger, nil
}
