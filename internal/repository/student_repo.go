package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/helpers"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var cacheTTL int = 5 // mins

type StudentRepository struct {
	collection *mongo.Collection
	cache      *cache.Cache
}

func NewStudentRepo(ctx context.Context, database *mongo.Database, c *cache.Cache) *StudentRepository {
	return &StudentRepository{
		collection: database.Collection("students"),
		cache:      c,
	}
}

func (r *StudentRepository) CreateStudent(ctx context.Context, student *models.Student, password string) (string, error) {
	if student == nil {
		return "", fmt.Errorf("nil student is provided")
	}

	// validate password
	if err := helpers.ValidatePassword(password); err != nil {
		return "", fmt.Errorf("password validate error : %w", err)
	}

	// hash password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password : %w", err)
	}

	timeNow := time.Now()

	// set student
	student.ID = primitive.NewObjectID()
	student.Role = "student"
	student.CreatedAt = timeNow
	student.UpdatedAt = timeNow
	student.IsActive = true
	student.PasswordHash = hashedPassword
	// init an empty slice
	student.RequiredExams = []primitive.ObjectID{}
	student.CompletedExams = []models.CompletedExam{}

	// add to mongo
	result, err := r.collection.InsertOne(ctx, student)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", fmt.Errorf("student with this ID is already exsists")
		}
		return "", err
	}

	studentID := result.InsertedID.(primitive.ObjectID).Hex()

	// set to cache
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:%s", student.StudentID)
		if err := r.cache.Set(cacheKey, student, time.Duration(cacheTTL)*time.Minute); err != nil {
			utils.LogErrorWithLevel("warn",
				utils.DragonflyFailedToWriteCache.Type,
				utils.DragonflyFailedToWriteCache.Code,
				utils.DragonflyFailedToWriteCache.Msg,
				err,
			)
		}
	}

	return studentID, nil
}

func (r *StudentRepository) fetchStudentFromDB(ctx context.Context, searchType, searchValue string) (*models.Student, error) {
	var student models.Student
	err := r.collection.FindOne(ctx, bson.M{searchType: searchValue}).Decode(&student)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) GetStudentByEmail(ctx context.Context, email string) (*models.Student, error) {
	var student models.Student

	cacheKey := fmt.Sprintf("user:email:%s", email)

	if r.cache == nil {
		return r.fetchStudentFromDB(ctx, "email", email)
	}

	err := r.cache.GetFromCacheOrFetchDB(
		ctx,
		cacheKey,
		&student,
		func() (any, error) {
			return r.fetchStudentFromDB(ctx, "email", email)
		},
		time.Duration(cacheTTL)*time.Minute,
	)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) GetStudentByID(ctx context.Context, studentID string) (*models.Student, error) {
	var student models.Student

	cacheKey := fmt.Sprintf("user:%s", studentID)

	if r.cache == nil {
		return r.fetchStudentFromDB(ctx, "student_id", studentID)
	}

	err := r.cache.GetFromCacheOrFetchDB(
		ctx,
		cacheKey,
		&student,
		func() (any, error) {
			return r.fetchStudentFromDB(ctx, "student_id", studentID)
		},
		time.Duration(cacheTTL)*time.Minute,
	)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) GetStudentByIDFromBD(ctx context.Context, studentID string) (*models.Student, error) {
	var student *models.Student
	var err error

	student, err = r.fetchStudentFromDB(ctx, "student_id", studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student by id from the database : %w", err)
	}

	return student, nil
}

func (r *StudentRepository) VerifyPassword(ctx context.Context, studentID, plainPassword string) (*models.Student, error) {
	student, err := r.GetStudentByIDFromBD(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("student not found: %w", err)
	}

	if !student.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	hashedPassword := student.PasswordHash

	if student.PasswordHash == "" {
		return nil, fmt.Errorf("password not set for this account")
	}

	// Check the password
	err = helpers.CheckWithHashedPassword(plainPassword, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return student, nil
}
