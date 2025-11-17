// Package controllers contains handler functions for the web application.
package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/dto/request"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Controllers struct {
	validator   *validator.Validate
	StudentRepo repository.StudentRepository
	cache       cache.Cache
}

var ctrl *Controllers

func NewControllers(valid *validator.Validate, studentRepo repository.StudentRepository, cache cache.Cache) {
	ctrl = &Controllers{
		valid,
		studentRepo,
		cache,
	}
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var createStudentRequest request.CreateStudentRequest

		// Get user input
		if err := ctx.BindJSON(&createStudentRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "error_details": err.Error()})
			return
		}

		// Validate user input
		if validationErr := ctrl.validator.Struct(createStudentRequest); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "error_details": validationErr.Error()})
			return
		}

		// TODO:
		// existingStudent, _ := ctrl.repo.GetUserByStudentID(c, createStudentRequest.StudentID)
		// if existingStudent != nil {
		// 	ctx.JSON(http.StatusConflict, gin.H{
		// 		"error": "Student with this ID already exists",
		// 	})
		// 	return
		// }

		// create student struct
		student := &models.Student{
			FirstName:  createStudentRequest.FirstName,
			LastName:   createStudentRequest.LastName,
			Email:      createStudentRequest.Email,
			StudentID:  createStudentRequest.StudentID,
			Department: createStudentRequest.Department,
		}

		StudentID, err := ctrl.StudentRepo.CreateStudent(c, student, createStudentRequest.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":         "Failed to create user account",
				"error_details": err.Error(),
			})
			return

		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "user has been created succussfuly", "student_id": StudentID})
	}
}
