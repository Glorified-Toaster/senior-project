package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/helpers"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	StudentColl    *mongo.Collection
	InstructorColl *mongo.Collection
	Cache          *cache.Cache
}

func NewStudentRepo(ctx context.Context, database *mongo.Database, c *cache.Cache) *UserRepository {
	return &UserRepository{
		StudentColl:    database.Collection("students"),
		InstructorColl: database.Collection("Instructors"),
		Cache:          c,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user any, password string) (string, error) {
	if user == nil {
		return "", fmt.Errorf("nil user is provided")
	}

	// Password hashing
	if err := helpers.ValidatePassword(password); err != nil {
		return "", fmt.Errorf("password validate error: %w", err)
	}

	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	timeNow := time.Now()
	var collection *mongo.Collection
	var idToReturn string
	var cacheKey string
	var documentToInsert any // ← Store what we'll actually insert

	// Identify the struct and set specific fields
	switch v := user.(type) {
	case *models.Student:
		v.ID = primitive.NewObjectID()
		v.Role = "student"
		v.PasswordHash = hashedPassword
		v.CreatedAt, v.UpdatedAt = timeNow, timeNow
		v.RequiredExams = []primitive.ObjectID{}
		collection = r.StudentColl
		idToReturn = v.ID.Hex()
		cacheKey = fmt.Sprintf("user:student_id:%s", v.StudentID)
		documentToInsert = v // ← Use the modified pointer

	case *models.Instructor:
		v.ID = primitive.NewObjectID()
		v.Role = "instructor"
		v.PasswordHash = hashedPassword
		v.CreatedAt, v.UpdatedAt = timeNow, timeNow
		collection = r.InstructorColl
		idToReturn = v.ID.Hex()
		cacheKey = fmt.Sprintf("user:instructor:instructor_id:%s", v.InstructorID)
		documentToInsert = v // ← Use the modified pointer

	default:
		return "", fmt.Errorf("unsupported type: %T", user)
	}

	// Database Logic - insert the MODIFIED document
	if _, err := collection.InsertOne(ctx, documentToInsert); err != nil {
		return "", err
	}

	// Cache Logic
	if r.Cache != nil {
		_ = r.Cache.Set(cacheKey, documentToInsert, time.Duration(cacheTTL)*time.Minute)
	}

	return idToReturn, nil
}

func (r *UserRepository) fetchFromDB(ctx context.Context, searchType, searchValue, role string) (models.User, error) {
	switch role {
	case "student":
		var student models.Student
		err := r.StudentColl.FindOne(ctx, bson.M{searchType: searchValue}).Decode(&student)
		if err != nil {
			return nil, err
		}
		return &student, nil

	case "instructor":
		var instructor models.Instructor
		err := r.InstructorColl.FindOne(ctx, bson.M{searchType: searchValue}).Decode(&instructor)
		if err != nil {
			return nil, err
		}
		return &instructor, nil

	default:
		return nil, fmt.Errorf("invalid role provided: %s", role)
	}
}

func (r *UserRepository) GetBy(ctx context.Context, searchType, searchValue string) (models.User, error) {
	cacheKey := fmt.Sprintf("user:%s:%s", searchType, searchValue)

	fetcher := func() (any, error) {
		student, err := r.fetchFromDB(ctx, searchType, searchValue, "student")
		if err == nil {
			return student, nil
		}

		instructor, err := r.fetchFromDB(ctx, searchType, searchValue, "instructor")
		if err == nil {
			return instructor, nil
		}

		return nil, fmt.Errorf("user not found with %s: %s", searchType, searchValue)
	}

	if r.Cache == nil {
		res, err := fetcher()
		if err != nil {
			return nil, err
		}
		return res.(models.User), nil
	}

	var foundUser any
	err := r.Cache.GetFromCacheOrFetchDB(
		ctx,
		cacheKey,
		&foundUser,
		fetcher,
		time.Duration(cacheTTL)*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	switch v := foundUser.(type) {
	case *models.Student:
		return v, nil
	case *models.Instructor:
		return v, nil
	case models.Student:
		// Value type from cache
		return &v, nil
	case models.Instructor:
		// Value type from cache
		return &v, nil
	case map[string]any:

		role, ok := v["role"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to determine user role from cache")
		}

		switch role {
		case "student":
			var student models.Student
			if err := mapToStruct(v, &student); err != nil {
				return nil, fmt.Errorf("failed to convert cached student: %w", err)
			}
			return &student, nil
		case "instructor":
			var instructor models.Instructor
			if err := mapToStruct(v, &instructor); err != nil {
				return nil, fmt.Errorf("failed to convert cached instructor: %w", err)
			}
			return &instructor, nil
		default:
			return nil, fmt.Errorf("unknown role in cache: %s", role)
		}
	default:
		return nil, fmt.Errorf("unexpected type from cache: %T", foundUser)
	}
}

func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func mapToStruct(m map[string]any, dest any) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *UserRepository) VerifyPassword(ctx context.Context, userID, plainPassword, role string) (models.User, error) {
	// Determine search type based on role
	searchType := "student_id"
	if role == "instructor" {
		searchType = "instructor_id"
	}

	// Fetch the user by appropriate ID
	user, err := r.GetBy(ctx, searchType, userID)
	if err != nil {
		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	hashedPassword := user.GetPassword()
	if hashedPassword == "" {
		return nil, fmt.Errorf("password not set for this account")
	}

	// Check the password
	err = helpers.CheckWithHashedPassword(plainPassword, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	return user, nil // Return the interface, not just student
}
