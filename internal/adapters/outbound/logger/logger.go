package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Logger {
	return &Logger{log: log}
}

// LogErrorWithLevel : log error and select the level of that error
func (logger *Logger) LogErrorWithLevel(level, errorType, errorCode, msg string, err error, fields ...zap.Field) {
	// must be passed fields
	mandatoryFields := []zap.Field{
		zap.String("error_code", errorCode),
		zap.String("error_msg", msg),
		zap.Error(err),
	}

	// append all of the additional fields
	allFields := append(mandatoryFields, fields...)

	switch level {

	case "fatal":
		logger.log.Fatal(errorType, allFields...)

	case "error":
		logger.log.Error(errorType, allFields...)

	case "warn":
		logger.log.Warn(errorType, allFields...)

	case "panic":
		logger.log.Panic(errorType, allFields...)

	default:
		logger.log.Error(errorType, allFields...)

	}
}

// LogInfo : log Info (very useful comment i guess...)
func (logger *Logger) LogInfo(infoType, msg string, fields ...zap.Field) {
	mandatoryFields := []zap.Field{
		zap.String("info_msg", msg),
	}

	allFields := append(mandatoryFields, fields...)
	logger.log.Info(infoType, allFields...)
}
