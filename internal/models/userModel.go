package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName   *string            `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName    *string            `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Password    *string            `bson:"password" json:"password" validate:"required,min=6"`
	Email       *string            `bson:"email" json:"email" validate:"email,required"`
	Phone       *string            `bson:"phone" json:"phone" validate:"required"`
	Token       *string            `bson:"token,omitempty" json:"token,omitempty"`
	Role        *string            `bson:"role" json:"role" validate:"required,eq=TEACHER|eq=USER"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UserID      string             `bson:"user_id" json:"user_id"`
	AccessToken string             `bson:"access_token" json:"access_token"`
}
