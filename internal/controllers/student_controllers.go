package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/dto/request"
	"github.com/Glorified-Toaster/senior-project/internal/dto/response"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/templates"
	"github.com/Glorified-Toaster/senior-project/internal/templates/components"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (ctrl *Controllers) StudentLoginPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		csrfToken := csrf.GetToken(ctx)

		render := utils.NewRender(ctx, http.StatusOK, templates.StudentLoginPage(csrfToken))
		ctx.Render(http.StatusOK, render)
	}
}

func (ctrl *Controllers) Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var createStudentRequest request.CreateStudentRequest

		// Get user input
		if err := ctx.BindJSON(&createStudentRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":         "invalid request format",
				"error_details": err.Error(),
			})
			return
		}

		// Validate user input
		if validationErr := ctrl.validator.Struct(createStudentRequest); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":         "validation failed",
				"error_details": validationErr.Error(),
			})
			return
		}

		// Check if student ID already exists
		existingStudent, err := ctrl.UserRepo.GetBy(c, "student_id", createStudentRequest.StudentID)
		if err == nil && existingStudent != nil {
			// User found - already exists
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Student with this ID already exists",
			})
			return
		}

		// Check if email already exists (important for login!)
		existingEmail, err := ctrl.UserRepo.GetBy(c, "email", createStudentRequest.Email)
		if err == nil && existingEmail != nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Email address is already registered",
			})
			return
		}

		// Create student struct
		student := &models.Student{
			FirstName:  createStudentRequest.FirstName,
			LastName:   createStudentRequest.LastName,
			Email:      createStudentRequest.Email,
			StudentID:  createStudentRequest.StudentID,
			Department: createStudentRequest.Department,
			IsActive:   true,
		}

		studentID, err := ctrl.UserRepo.CreateUser(c, student, createStudentRequest.Password)
		if err != nil {
			// Check if it's a duplicate key error from MongoDB
			if mongo.IsDuplicateKeyError(err) {
				ctx.JSON(http.StatusConflict, gin.H{
					"error": "Student ID or email already exists",
				})
				return
			}

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
			utils.LogErrorWithLevel("error", "HTTP_SERVER", "JWT_GEN_FAILED_ERROR",
				"failed to generate JWT token after signup", err)
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

		user, err := ctrl.UserRepo.GetBy(ctx.Request.Context(), "student_id", studentID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "student not found",
			})
			return
		}

		student, ok := user.(*models.Student)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID belongs to an instructor, not a student"})
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
		c, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
		defer cancel()

		var loginReq request.LoginRequest

		if err := ctx.ShouldBind(&loginReq); err != nil {
			render := utils.NewRender(ctx, http.StatusBadRequest, components.ErrorToast("Invalid form data"))
			ctx.Render(http.StatusBadRequest, render)
			return
		}

		// Validate required fields
		if loginReq.UserID == "" || loginReq.Password == "" {
			render := utils.NewRender(ctx, http.StatusBadRequest, components.ErrorToast("User ID and password are required"))
			ctx.Render(http.StatusBadRequest, render)
			return
		}

		user, err := ctrl.UserRepo.VerifyPassword(c, loginReq.UserID, loginReq.Password, "student")
		if err != nil {
			utils.LogInfo("HTTP_SERVER", "Login Failed",
				zap.String("user_id", loginReq.UserID),
				zap.String("error", err.Error()))

			render := utils.NewRender(ctx, http.StatusOK, components.ErrorToast("Login Failed! Check credentials."))
			ctx.Render(http.StatusOK, render)
			return
		}

		// Type assertion to get student
		student, ok := user.(*models.Student)
		if !ok {
			render := utils.NewRender(ctx, http.StatusOK, components.ErrorToast("Invalid user type"))
			ctx.Render(http.StatusOK, render)
			return
		}

		// Check if the account is active
		if !student.IsActive {
			render := utils.NewRender(ctx, http.StatusOK, components.ErrorToast("Account is not active"))
			ctx.Render(http.StatusOK, render)
			return
		}

		now := time.Now()
		student.LastLogin = &now

		// Generate JWT token
		additionalClaims := map[string]any{
			"first_name": student.FirstName,
			"last_name":  student.LastName,
			"department": student.Department,
			"student_id": student.StudentID,
			"is_active":  student.IsActive,
			"user_id":    student.ID.Hex(),
		}

		token, err := ctrl.jwtAuth.GenerateToken(student.Email, student.ID.Hex(), "student", additionalClaims)
		if err != nil {
			utils.LogErrorWithLevel("error", "HTTP_SERVER", "JWT_GEN_FAILED_ERROR",
				"failed to generate JWT token after login", err)

			render := utils.NewRender(ctx, http.StatusOK, components.ErrorToast("Authentication failed"))
			ctx.Render(http.StatusOK, render)
			return
		}

		// Set cookie for web requests
		ctx.SetCookie("auth_token", token, 86400, "/", "", true, true) // Secure: true in production

		// For HTML requests, redirect or show success
		render := utils.NewRender(ctx, http.StatusOK,
			components.SuccessToast("Login successful! Redirecting..."))
		ctx.Render(http.StatusOK, render)
	}
}
