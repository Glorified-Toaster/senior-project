package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/adapters/outbound/logger"
	"uot-exam/internal/application"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/internal/utils/random"
	"uot-exam/web/templates/components/toast"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/components"
	"uot-exam/web/templates/render"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type UserHandler struct {
	App         *application.Application
	validate    *validator.Validate
	jwt         *helpers.JWTAuth
	viperConfig *config.Config
	logger      *logger.Logger
}

func NewUserHandler(App *application.Application, validate *validator.Validate, jwt *helpers.JWTAuth, viperConfig *config.Config, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		App:         App,
		validate:    validate,
		jwt:         jwt,
		viperConfig: viperConfig,
		logger:      logger,
	}
}

func (h *UserHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fullname := ctx.PostForm("full_name")
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		role := ctx.PostForm("role")

		if helpers.IsTrimmedEmpty(fullname) ||
			helpers.IsTrimmedEmpty(username) ||
			helpers.IsTrimmedEmpty(password) ||
			helpers.IsTrimmedEmpty(role) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create User Failed", "All fields are required.", toast.VariantError)
			return
		}

		_, err := h.App.CreateUser(ctx, ports.CreateUserParams{
			FullName: fullname,
			Username: username,
			Password: password,
			Role:     domain.UserRole(role),
			IsActive: true,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create User Failed", "Error creating user: "+err.Error(), toast.VariantError)
			return
		}

		ctx.Header("HX-Reswap", "none")
		helpers.Toast(ctx, "Create User Success", "User created successfully", toast.VariantSuccess)

		// Fetch updated first page of users to refresh the table
		users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{
			Limit:  12,
			Offset: 0,
		})
		if err != nil {
			users = []domain.User{}
		}

		totalCount, err := h.App.CountUsers(ctx)
		if err != nil {
			totalCount = 0
		}

		render.Render(ctx, components.UserTableContainer(components.UserTableProps{
			Users:         users,
			Title:         "All Users",
			ID:            "users-table",
			TotalCount:    totalCount,
			Limit:         12,
			Offset:        0,
			BaseURL:       "/admin/dashboard/users",
			Search:        true,
			SearchAPI:     "/admin/users/search",
			AddUser:       true,
			DeletedButton: true,
		}))
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ports.LoginParams
		if err := ctx.ShouldBind(&req); err != nil {
			toast.Toast(toast.Props{
				Title:         "Login Failed",
				Description:   "Invalid request body",
				Variant:       toast.VariantError,
				Duration:      4000,
				ShowIndicator: true,
				Dismissible:   true,
				Icon:          true,
			}).Render(ctx.Request.Context(), ctx.Writer)
			return
		}

		if err := h.validate.Struct(req); err != nil {
			toast.Toast(toast.Props{
				Title:         "Login Failed",
				Description:   "Invalid request body validation",
				Variant:       toast.VariantError,
				Duration:      4000,
				ShowIndicator: true,
				Dismissible:   true,
				Icon:          true,
			}).Render(ctx.Request.Context(), ctx.Writer)
			return
		}

		user, err := h.App.Login(ctx, req)
		if err != nil {
			toast.Toast(toast.Props{
				Title:         "Login Failed",
				Description:   "Invalid username or password",
				Variant:       toast.VariantError,
				Duration:      4000,
				ShowIndicator: true,
				Dismissible:   true,
				Icon:          true,
			}).Render(ctx.Request.Context(), ctx.Writer)
			return
		}

		token, err := h.jwt.GenerateToken(user)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":        "User logged in successfully. Please login to get access token.",
				"student_id": user.ID,
				"warning":    "Token generation failed - please login manually",
			})
			return
		}

		ctx.SetCookie("auth_token", token, 60*60*24, "/", "", false, true)
		ctx.Header("Content-Type", "text/html; charset=utf-8")

		switch user.Role {
		case domain.RoleAdmin:
			ctx.Header("HX-Redirect", "/admin/dashboard")
		case domain.RoleStudent:
			ctx.Header("HX-Redirect", "/student/dashboard")
		default:
			ctx.Header("HX-Redirect", "/login")
		}
	}
}

func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "INVALID_REQUEST", "INVALID_REQUEST", "Invalid request body", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := h.App.GetUserByID(ctx, id)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to get user by ID", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get user by ID"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (h *UserHandler) GetUserByUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		if username == "" {
			h.logger.LogErrorWithLevel("warn", "INVALID_REQUEST", "INVALID_REQUEST", "Invalid request body", nil)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := h.App.GetUserByUsername(ctx, username)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to get user by username", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get user by username"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (h *UserHandler) SearchUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		search := strings.TrimSpace(ctx.PostForm("search"))
		if search == "" {
			search = strings.TrimSpace(ctx.Query("search"))
		}
		limitStr := ctx.Query("limit")
		offsetStr := ctx.Query("offset")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 12
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		// If search is empty, return the first page with pagination
		if search == "" {
			users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{
				Limit:  int32(limit),
				Offset: int32(offset),
			})
			if err != nil {
				users = []domain.User{}
			}
			totalCount, _ := h.App.CountUsers(ctx)
			ctx.Header("Content-Type", "text/html")
			render.Render(ctx, components.UserTableContainer(components.UserTableProps{
				Users:      users,
				TotalCount: totalCount,
				Limit:      int32(limit),
				Offset:     int32(offset),
				BaseURL:    "/admin/dashboard/users",
			}))
			return
		}

		users, err := h.App.SearchUsers(ctx, ports.SearchUsersParams{
			Search: search,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "SEARCH_FAILED", "Failed to search users", err)
			render.Render(ctx, components.UserTableRows([]domain.User{}, false, false))
			return
		}

		totalCount, err := h.App.CountSearchUsers(ctx, search)
		if err != nil {
			totalCount = 0
		}

		if users == nil {
			users = []domain.User{}
		}
		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.UserTableContainer(components.UserTableProps{
			Users:         users,
			TotalCount:    totalCount,
			Limit:         int32(limit),
			Offset:        int32(offset),
			BaseURL:       "/admin/users/search?search=" + url.QueryEscape(search),
			Search:        true,
			SearchAPI:     "/admin/users/search",
			ShowAllButton: false,
		}))
	}
}

func (h *UserHandler) ListAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{Limit: 100, Offset: 0})
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to list all users", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to list all users"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}

func (h *UserHandler) TestPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitStr := ctx.Query("limit")
		offsetStr := ctx.Query("offset")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 12
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to list all users", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to list all users"})
			return
		}

		totalCount, err := h.App.CountUsers(ctx)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to count users", err)
			totalCount = 0
		}

		if ctx.GetHeader("HX-Request") == "true" {
			render.Render(ctx, pages.UserTableContainer(users, totalCount, int32(limit), int32(offset)))
			return
		}

		render.Render(ctx, pages.TestPage(users, totalCount, int32(limit), int32(offset)))
	}
}

func (h *UserHandler) SoftDeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "INVALID_REQUEST", "INVALID_REQUEST", "Invalid request body", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid request body"})
			return
		}
		if err = h.App.SoftDeleteUser(ctx, id); err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to delete user", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete user"})
			return
		}

		limit := 12
		offset := 0

		if limitStr := ctx.Query("limit"); limitStr != "" {
			if parsed, err := strconv.Atoi(limitStr); err == nil {
				limit = parsed
			}
		}
		if offsetStr := ctx.Query("offset"); offsetStr != "" {
			if parsed, err := strconv.Atoi(offsetStr); err == nil {
				offset = parsed
			}
		}

		if currentURL := ctx.Request.Header.Get("HX-Current-URL"); currentURL != "" {
			if u, parseErr := url.Parse(currentURL); parseErr == nil {
				q := u.Query()
				if limitStr := q.Get("limit"); limitStr != "" {
					if parsed, err := strconv.Atoi(limitStr); err == nil {
						limit = parsed
					}
				}
				if offsetStr := q.Get("offset"); offsetStr != "" {
					if parsed, err := strconv.Atoi(offsetStr); err == nil {
						offset = parsed
					}
				}
			}
		}

		users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to list users after delete", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to list users"})
			return
		}

		totalCount, err := h.App.CountUsers(ctx)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to count users after delete", err)
			totalCount = 0
		}

		props := components.UserTableProps{
			Title:      "All Users",
			Users:      users,
			TotalCount: totalCount,
			Limit:      int32(limit),
			Offset:     int32(offset),
			BaseURL:    "/admin/dashboard/users",
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.UserTableContainerWithDeleteToast(props))
	}
}

func (h *UserHandler) SearchDeletedUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		search := strings.TrimSpace(ctx.PostForm("search"))
		if search == "" {
			search = strings.TrimSpace(ctx.Query("search"))
		}
		limitStr := ctx.Query("limit")
		offsetStr := ctx.Query("offset")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		// If search is empty, return the first page with pagination
		if search == "" {
			users, err := h.App.ListDeletedUsers(ctx, ports.ListDeletedUsersParams{
				Limit:  int32(limit),
				Offset: int32(offset),
			})
			if err != nil {
				users = []domain.User{}
			}
			totalCount, _ := h.App.CountDeletedUsers(ctx)
			ctx.Header("Content-Type", "text/html")
			render.Render(ctx, components.UserTableContainer(components.UserTableProps{
				ID:            "deleted-users-table-container",
				Users:         users,
				TotalCount:    totalCount,
				Limit:         int32(limit),
				Offset:        int32(offset),
				BaseURL:       "/admin/dashboard/users/deleted",
				RestoreButton: true,
			}))
			return
		}

		users, err := h.App.SearchDeletedUsers(ctx, ports.SearchUsersParams{
			Search: search,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		// Get total count for search results
		totalCount, err := h.App.CountSearchDeletedUsers(ctx, search)
		if err != nil {
			totalCount = 0
		}

		if users == nil {
			users = []domain.User{}
		}
		render.Render(ctx, components.UserTableContainer(components.UserTableProps{
			ID:            "deleted-users-table-container",
			Users:         users,
			TotalCount:    totalCount,
			Limit:         int32(limit),
			Offset:        int32(offset),
			BaseURL:       "/admin/dashboard/users/deleted/search?search=" + url.QueryEscape(search),
			Search:        true,
			SearchAPI:     "/admin/users/search-deleted",
			ShowAllButton: false,
			RestoreButton: true,
		}))
	}
}

func (h *UserHandler) GenerateRandomUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		suffix, err := random.String(8)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "INTERNAL_ERROR", "INTERNAL_ERROR", "Failed to generate random username", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.String(http.StatusOK, "user_"+suffix)
	}
}

func (h *UserHandler) GenerateRandomPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pass, err := random.Password(12)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "INTERNAL_ERROR", "INTERNAL_ERROR", "Failed to generate random password", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.String(http.StatusOK, pass)
	}
}

func (h *UserHandler) ToggleUserActive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "INVALID_REQUEST", "INVALID_REQUEST", "Invalid user ID", err)
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Toggle User Failed", "Invalid user ID", toast.VariantError)
			return
		}

		user, err := h.App.GetUserByID(ctx, id)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to get user", err)
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Toggle User Failed", "Failed to get user", toast.VariantError)
			return
		}

		if user.IsActive {
			err = h.App.DisableUser(ctx, id)
		} else {
			err = h.App.EnableUser(ctx, id)
		}

		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to toggle user status", err)
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Toggle User Failed", "Failed to toggle user status: "+err.Error(), toast.VariantError)
			return
		}

		// Reload users list
		limit := int32(12)
		offset := int32(0)

		users, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			users = []domain.User{}
		}

		totalCount, err := h.App.CountUsers(ctx)
		if err != nil {
			totalCount = 0
		}

		action := "disabled"
		if !user.IsActive {
			action = "enabled"
		}

		ctx.Header("Content-Type", "text/html")
		helpers.Toast(ctx, "Success", "User "+action+" successfully", toast.VariantSuccess)
		render.Render(ctx, components.UserTableContainer(components.UserTableProps{
			Users:         users,
			Title:         "All Users",
			ID:            "users-table",
			TotalCount:    totalCount,
			Limit:         limit,
			Offset:        offset,
			BaseURL:       "/admin/dashboard/users",
			Search:        true,
			SearchAPI:     "/admin/users/search",
			AddUser:       true,
			DeletedButton: true,
		}))
	}
}

func (h *UserHandler) RestoreUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Restore User Failed", "Invalid user ID", toast.VariantError)
			return
		}

		if err = h.App.RestoreUser(ctx, id); err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Restore User Failed", "Failed to restore user: "+err.Error(), toast.VariantError)
			return
		}

		// Reload deleted users list
		deletedUsers, err := h.App.ListDeletedUsers(ctx, ports.ListDeletedUsersParams{
			Limit:  12,
			Offset: 0,
		})
		if err != nil {
			deletedUsers = []domain.User{}
		}

		totalCount, _ := h.App.CountDeletedUsers(ctx)

		ctx.Header("Content-Type", "text/html")
		helpers.Toast(ctx, "User Restored", "The user has been restored successfully.", toast.VariantSuccess)
		render.Render(ctx, components.UserTableContainer(components.UserTableProps{
			ID:            "deleted-users-table-container",
			Users:         deletedUsers,
			TotalCount:    totalCount,
			Limit:         12,
			Offset:        0,
			BaseURL:       "/admin/dashboard/users/deleted",
			Search:        true,
			SearchAPI:     "/admin/dashboard/users/deleted/search",
			RestoreButton: true,
		}))
	}
}
