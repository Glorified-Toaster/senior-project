// Package repository
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	ctx        context.Context
	collection *mongo.Collection
	cache      *cache.Cache
}

// NewUserRepo : constructor
func NewUserRepo(ctx context.Context, database *mongo.Database, c *cache.Cache) *UserRepository {
	return &UserRepository{
		ctx:        ctx,
		collection: database.Collection("students"),
		cache:      c,
	}
}

// GetUserByID : implements cache-aside pattern
func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	cacheKey := fmt.Sprintf("user:%s", userID)

	// If cache isn't configured, read directly from DB
	if r.cache == nil {
		data, err := r.fetchUserFromDB(r.ctx, userID)
		if err != nil {
			return nil, err
		}
		if user, ok := data.(models.User); ok {
			return &user, nil
		}
		// defensive: try to convert if underlying type differs
		return nil, fmt.Errorf("unexpected user type returned from DB")
	}

	// Cache-aside: try cache, then DB
	err := r.cache.GetFromCacheOrFetchDB(
		r.ctx,
		cacheKey,
		&user,
		func() (any, error) {
			// This function is called on cache miss
			return r.fetchUserFromDB(r.ctx, userID)
		},
		15*time.Minute, // Cache for 15 minutes
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) fetchUserFromDB(ctx context.Context, userID string) (any, error) {
	var user models.User

	// if userID is a hex ObjectID string, you can query _id
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
		if err == nil {
			return user, nil
		}
		// if not found by _id, continue to try user_id below
	}

	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID string, updates bson.M) error {
	filter := bson.M{"user_id": userID}

	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		filter = bson.M{"_id": oid}
	}

	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": updates})
	if err != nil {
		return err
	}

	// invalidate cache if present
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:%s", userID)
		if derr := r.cache.Delete(cacheKey); derr != nil {
			utils.LogErrorWithLevel("error",
				utils.DragonflyFailedToDeleteCache.Type,
				utils.DragonflyFailedToDeleteCache.Code,
				utils.DragonflyFailedToDeleteCache.Msg,
				derr)
		}
	}
	return nil
}

// DeleteUser deletes a user document and invalidates cache
func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}
	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		filter = bson.M{"_id": oid}
	}

	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// invalidate cache if present
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:%s", userID)
		if derr := r.cache.Delete(cacheKey); derr != nil {
			utils.LogErrorWithLevel("warn",
				utils.DragonflyFailedToDeleteCache.Type,
				utils.DragonflyFailedToDeleteCache.Code,
				"failed to delete cache key after user delete",
				derr)
		}
	}

	return nil
}

// CreateUser inserts a new user document and optionally writes it to cache.
// It sets CreatedAt, UpdatedAt and a user_id if missing.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("nil user provided")
	}

	now := time.Now()
	user.UpdatedAt = now
	user.CreatedAt = now

	// ensure an ObjectID
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	// ensure a user_id string (use hex of ObjectID if not set)
	if user.UserID == "" {
		user.UserID = user.ID.Hex()
	}

	// insert into MongoDB
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	// attempt to write to cache (best-effort)
	if r.cache != nil {
		cacheKey := fmt.Sprintf("user:%s", user.UserID)
		if serr := r.cache.Set(cacheKey, user, 15*time.Minute); serr != nil {
			utils.LogErrorWithLevel("warn",
				utils.DragonflyFailedToWriteCache.Type,
				utils.DragonflyFailedToWriteCache.Code,
				utils.DragonflyFailedToWriteCache.Msg,
				serr)
		}
	}

	// return the inserted id (prefer user_id)
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return user.UserID, nil
}
