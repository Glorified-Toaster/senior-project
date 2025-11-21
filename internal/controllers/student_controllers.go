package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/dto/request"
	"github.com/Glorified-Toaster/senior-project/internal/dto/response"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}

func (ctrl *Controllers) Signup() gin.HandlerFunc {
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

		existingStudent, _ := ctrl.StudentRepo.GetStudentByID(c, createStudentRequest.StudentID)
		if existingStudent != nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Student with this ID already exists",
			})
			return
		}

		// create student struct
		student := &models.Student{
			FirstName:  createStudentRequest.FirstName,
			LastName:   createStudentRequest.LastName,
			Email:      createStudentRequest.Email,
			StudentID:  createStudentRequest.StudentID,
			Department: createStudentRequest.Department,
			IsActive:   true,
		}

		studentID, err := ctrl.StudentRepo.CreateStudent(c, student, createStudentRequest.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":         "Failed to create user account",
				"error_details": err.Error(),
			})
			return

		}

		additionalClaims := map[string]any{
			"first_name": student.FirstName,
			"last_name":  student.LastName,
			"department": student.Department,
			"student_id": student.StudentID,
			"is_active":  student.IsActive,
		}

		token, err := ctrl.jwtAuth.GenerateToken(student.Email, studentID, "student", additionalClaims)
		if err != nil {
			utils.LogErrorWithLevel("error", "HTTP_SERVER", "JWT_GEN_FAILED_ERROR", "failed to generate JWT token after signup", err)

			ctx.JSON(http.StatusOK, gin.H{
				"msg":        "User created successfully. Please login to get access token.",
				"student_id": studentID,
				"warning":    "Token generation failed - please login manually",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"msg":          "User created successfully",
			"student_id":   studentID,
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   24 * 60 * 60,
			"user": gin.H{
				"id":         studentID,
				"first_name": student.FirstName,
				"last_name":  student.LastName,
				"email":      student.Email,
				"student_id": student.StudentID,
				"department": student.Department,
				"role":       "student",
			},
		})
	}
}

func (ctrl *Controllers) GetStudentByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		studentID := ctx.Param("id")

		if studentID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "student ID is required",
			})
			return
		}

		student, err := ctrl.StudentRepo.GetStudentByID(ctx.Request.Context(), studentID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "student not found",
			})
			return
		}

		studentResponse := response.StudentResponse{
			ID:         student.ID,
			FirstName:  student.FirstName,
			LastName:   student.LastName,
			Role:       student.Role,
			Department: student.Department,
			StudentID:  student.StudentID,
			Email:      student.Email,
			IsActive:   student.IsActive,
			LastLogin:  student.LastLogin,
			CreatedAt:  student.CreatedAt,
			UpdatedAt:  student.UpdatedAt,
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Student retrieved successfully",
			"data":    studentResponse,
		})
	}
}

func (ctrl *Controllers) StudentLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var loginRequest *request.StudentLoginRequest

		if err := ctx.BindJSON(&loginRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		studentID := loginRequest.StudentID
		password := strings.TrimSpace(loginRequest.Password)

		student, err := ctrl.StudentRepo.VerifyPassword(c, studentID, password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", "error_msg": err.Error()})
			return
		}

		additionalClaims := map[string]any{
			"first_name": student.FirstName,
			"last_name":  student.LastName,
			"department": student.Department,
			"student_id": student.StudentID,
			"is_active":  student.IsActive,
		}

		token, err := ctrl.jwtAuth.GenerateToken(student.Email, student.StudentID, "student", additionalClaims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   86400,
			"user": gin.H{
				"id":         student.ID,
				"first_name": student.FirstName,
				"last_name":  student.LastName,
				"email":      student.Email,
				"student_id": student.StudentID,
				"department": student.Department,
				"role":       "student",
			},
		})
	}
}
