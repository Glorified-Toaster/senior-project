package application

import (
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/logger"
	"uot-exam/internal/ports"
)

type Application struct {
	userRepo  ports.UserRepository
	txManager ports.TransactionManager
	db        *database.PostgresAdapter
	log       *logger.Logger
}

func NewApplication(userRepo ports.UserRepository, txManager ports.TransactionManager, db *database.PostgresAdapter, log *logger.Logger) *Application {
	return &Application{userRepo: userRepo, txManager: txManager, db: db, log: log}
}
