package models

type User interface {
	GetID() string
	GetRole() string
	GetEmail() string
	GetPassword() string
	IsActiveUser() bool
}
