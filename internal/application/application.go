package application

import (
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/logger"
	"uot-exam/internal/ports"
)

type Application struct {
	userRepo     ports.UserRepository
	subjectRepo  ports.SubjectRepository
	examRepo     ports.ExamRepository
	questionRepo ports.QuestionRepository
	txManager    ports.TransactionManager
	localDisk    ports.LocalDiskAdapter
	db           *database.PostgresAdapter
	log          *logger.Logger
}

func NewApplication(userRepo ports.UserRepository,
	subjectRepo ports.SubjectRepository,
	examRepo ports.ExamRepository,
	questionRepo ports.QuestionRepository,
	localDisk ports.LocalDiskAdapter,
	txManager ports.TransactionManager,
	db *database.PostgresAdapter,
	log *logger.Logger) *Application {
	return &Application{
		userRepo:     userRepo,
		subjectRepo:  subjectRepo,
		examRepo:     examRepo,
		questionRepo: questionRepo,
		localDisk:    localDisk,
		txManager:    txManager,
		db:           db,
		log:          log,
	}
}

func (a *Application) DB() *database.PostgresAdapter {
	return a.db
}
