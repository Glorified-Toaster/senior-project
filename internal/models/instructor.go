package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instructor struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName    string             `bson:"first_name" json:"first_name" validate:"required,min=2,max=32"`
	LastName     string             `bson:"last_name" json:"last_name" validate:"required,min=2,max=32"`
	Role         string             `bson:"role" json:"role"`
	Department   string             `bson:"department,omitempty" json:"department,omitempty"`
	InstructorID string             `bson:"instructor_id" json:"instructor_id"`
	Email        string             `bson:"email" json:"email" validate:"email,required"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	LastLogin    *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func (i *Instructor) GetID() string       { return i.InstructorID }
func (i *Instructor) GetRole() string     { return "instructor" }
func (i *Instructor) GetEmail() string    { return i.Email }
func (i *Instructor) GetPassword() string { return i.PasswordHash }
func (i *Instructor) IsActiveUser() bool  { return i.IsActive }
