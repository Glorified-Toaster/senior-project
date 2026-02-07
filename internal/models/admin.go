package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Privileges struct {
	Name   string `bson:"name,omitempty"`
	Status bool   `bson:"status,omitempty"`
}

// type SuperAdmin struct {
// 	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	FullName     string             `bson:"fullname" json:"fullname" validate:"required,min=2,max=32"`
// 	UserName     string             `bson:"username" json:"username" validate:"required,min=2,max=16"`
// 	Role         string             `bson:"role" json:"role"`
// 	PasswordHash string             `bson:"password_hash" json:"-"`
// 	LastLogin    *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
// 	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
// 	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
// }

type Admin struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FullName      string             `bson:"fullname" json:"fullname" validate:"required,min=2,max=32"`
	UserName      string             `bson:"username" json:"username" validate:"required,min=2,max=16"`
	Role          string             `bson:"role" json:"role"`
	IsActive      bool               `bson:"is_active" json:"is_active"`
	PasswordHash  string             `bson:"password_hash" json:"-"`
	LastLogin     *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	HasPrivileges []Privileges       `bson:"has_privileges,omitempty" json:"has_privileges,omitempty"`
}
