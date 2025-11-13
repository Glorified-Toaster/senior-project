package utils

import (
	"sync"

	"github.com/Glorified-Toaster/senior-project/internal/config/logger"
	"go.uber.org/zap"
)

var (
	zlogger *zap.Logger
	once    sync.Once
)

// InitUtils : get logger (once)
func InitUtils() {
	once.Do(func() {
		zlogger = logger.GetLogger()
	})
}
