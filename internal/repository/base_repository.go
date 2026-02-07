package repository

import (
	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

var cacheTTL int = 5

type BaseRepository struct {
	Collection *mongo.Collection
	Cache      *cache.Cache
}
