package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"uot-exam/internal/adapters/outbound/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitZapLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg == nil || cfg.ZapLogger == nil {
		return nil, fmt.Errorf("invalid logger configuration")
	}

	logDir := cfg.ZapLogger.DirPath

	if err := ensureLogDirectory(logDir); err != nil {
		return nil, fmt.Errorf("failed to setup log directory: %w", err)
	}

	logFile := filepath.Join(logDir, cfg.ZapLogger.FileName)

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.Lumberjack.MaxSize,
		MaxBackups: cfg.Lumberjack.MaxBackups,
		MaxAge:     cfg.Lumberjack.MaxAge,
		Compress:   cfg.Lumberjack.Compress,
	}

	level, err := zapcore.ParseLevel(cfg.ZapLogger.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid zap level '%s': %w", cfg.ZapLogger.Level, err)
	}

	consoleEncoderConfig := getConsoleEncoderConfig(cfg.ZapLogger.Development)
	fileEncoderConfig := getFileEncoderConfig()

	// Create encoders
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	// Create cores
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(lumberjackLogger),
		level,
	)

	core := zapcore.NewTee(consoleCore, fileCore)

	zapLogger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	zapLogger.Info("Logger initialized successfully",
		zap.String("level", cfg.ZapLogger.Level),
		zap.String("log_file", logFile),
		zap.Bool("development", cfg.ZapLogger.Development),
	)

	return zapLogger, nil
}

func getConsoleEncoderConfig(isDevelopment bool) zapcore.EncoderConfig {
	var config zapcore.EncoderConfig

	if isDevelopment {
		config = zap.NewDevelopmentEncoderConfig()
	} else {
		config = zap.NewProductionEncoderConfig()
	}

	// Colors for terminal
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder // ✅ Colors for console
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder
	config.ConsoleSeparator = " | "

	return config
}

func getFileEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()

	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	config.MessageKey = "message"
	config.LevelKey = "level"
	config.TimeKey = "timestamp"
	config.CallerKey = "caller"
	config.StacktraceKey = "stacktrace"

	return config
}

func ensureLogDirectory(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to access log directory: %w", err)
	}
	return nil
}
