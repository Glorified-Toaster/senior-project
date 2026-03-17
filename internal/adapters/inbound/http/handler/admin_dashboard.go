package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/components"
	"uot-exam/web/templates/pages/admin_dashboard/page"
	"uot-exam/web/templates/render"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) AdminDashboardMainRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.App.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{Limit: 4, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		username, fullname := parseUsername(ctx)
		description := "something"
		exams := []domain.Exam{
			{
				ID:              uuid.New(),
				Title:           "Theory of computation",
				ExamID:          "TOC-001",
				Description:     &description,
				DurationMinutes: 60,
				TotalMarks:      100,
				Status:          domain.ExamStatusPublished,
				CreatedBy:       uuid.New(),
				StartTime:       time.Now(),
				EndTime:         time.Now().Add(1 * time.Hour),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              uuid.New(),
				Title:           "Theory of automatas",
				ExamID:          "TOA-001",
				Description:     &description,
				DurationMinutes: 60,
				TotalMarks:      100,
				Status:          domain.ExamStatusPublished,
				CreatedBy:       uuid.New(),
				StartTime:       time.Now(),
				EndTime:         time.Now().Add(1 * time.Hour),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              uuid.New(),
				Title:           "Mathematical Logic",
				ExamID:          "ML-001",
				Description:     &description,
				DurationMinutes: 60,
				TotalMarks:      100,
				Status:          domain.ExamStatusClosed,
				CreatedBy:       uuid.New(),
				StartTime:       time.Now(),
				EndTime:         time.Now().Add(1 * time.Hour),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              uuid.New(),
				Title:           "Data Structures",
				ExamID:          "DS-001",
				Description:     &description,
				DurationMinutes: 60,
				TotalMarks:      100,
				Status:          domain.ExamStatusDraft,
				CreatedBy:       uuid.New(),
				StartTime:       time.Now(),
				EndTime:         time.Now().Add(1 * time.Hour),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              uuid.New(),
				Title:           "Computer Graphics",
				ExamID:          "CG-001",
				Description:     &description,
				DurationMinutes: 60,
				TotalMarks:      100,
				Status:          domain.ExamStatusDraft,
				CreatedBy:       uuid.New(),
				StartTime:       time.Now(),
				EndTime:         time.Now().Add(1 * time.Hour),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}

		userCount, err := h.App.CountUsers(ctx)
		if err != nil {
			return
		}

		params := page.AdminDashboardParam{
			Users:           users,
			Username:        username,
			FullName:        fullname,
			Exams:           exams,
			TotalUsers:      fmt.Sprintf("%d", userCount),
			TotalUsersCount: userCount,
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AdminMainPage(params)))
	}
}

func (h *UserHandler) AdminLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		render.Render(ctx, pages.BasePage("Admin Login", page.LoginPage("")))
	}
}

func (h *UserHandler) UserPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitStr := ctx.DefaultQuery("limit", "10")
		offsetStr := ctx.DefaultQuery("offset", "0")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		users, err := h.App.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userCount, err := h.App.CountUsers(ctx)
		if err != nil {
			return
		}

		props := components.UserTableProps{
			Title:      "All Users",
			Users:      users,
			TotalCount: userCount,
			Limit:      int32(limit),
			Offset:     int32(offset),
			BaseURL:    "/admin/dashboard/users",
		}

		// Check if it's an HTMX request
		if ctx.GetHeader("HX-Request") != "" {
			render.Render(ctx, components.UserTableContainer(props))
			return
		}

		username, fullname := parseUsername(ctx)

		params := page.AdminDashboardParam{
			Users:           users,
			TotalUsers:      fmt.Sprintf("%d", userCount),
			TotalUsersCount: userCount,
			Username:        username,
			FullName:        fullname,
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AllUsers(params)))
	}
}

func (h *UserHandler) DeletedUsersPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitStr := ctx.DefaultQuery("limit", "10")
		offsetStr := ctx.DefaultQuery("offset", "0")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		deletedUsers, err := h.App.ListDeletedUsers(ctx.Request.Context(), ports.ListDeletedUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		deletedUserCount, err := h.App.CountDeletedUsers(ctx)
		if err != nil {
			return
		}

		// Check if it's an HTMX request
		if ctx.GetHeader("HX-Request") != "" {
			props := components.UserTableProps{
				Title:      "Deleted Users",
				Users:      deletedUsers,
				TotalCount: deletedUserCount,
				Limit:      int32(limit),
				Offset:     int32(offset),
				BaseURL:    "/admin/dashboard/users/deleted",
			}
			render.Render(ctx, components.UserTableContainer(props))
			return
		}

		username, fullname := parseUsername(ctx)

		params := page.DeletedUsersPageParams{
			DeletedUsers:    deletedUsers,
			TotalUsers:      fmt.Sprintf("%d", deletedUserCount),
			TotalUsersCount: deletedUserCount,
			Username:        username,
			FullName:        fullname,
			Title:           "Deleted Users",
		}
		render.Render(ctx, pages.BasePage("Deleted Users", page.DeletedUsersPage(params)))
	}
}

func parseUsername(ctx *gin.Context) (string, string) {
	usernameVal, exists := ctx.Get("username")
	if !exists {
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
		return "", ""
	}
	username := usernameVal.(string)
	fullnameVal, exists := ctx.Get("fullname")
	if !exists {
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
		fmt.Println("fullname not found")
		return "", ""
	}
	fullname := fullnameVal.(string)

	return username, fullname
}

func (h *UserHandler) AllExamsPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		render.Render(ctx, pages.BasePage("All Exams", page.AllExamsPage(page.AllExamsPageParam{})))
	}
}
