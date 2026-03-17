package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/components"
	"uot-exam/web/templates/pages/admin_dashboard/page"
	"uot-exam/web/templates/render"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"uot-exam/internal/adapters/inbound/http/helpers"
)

func (h *UserHandler) AdminDashboardMainRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.App.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{Limit: 4, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		username, fullname := parseUsername(ctx)
		exams, err := h.App.ListAllExams(ctx.Request.Context(), ports.ListAllExamsParams{Limit: 4, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
		limitStr := ctx.DefaultQuery("limit", "10")
		offsetStr := ctx.DefaultQuery("offset", "0")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		exams, err := h.App.ListAllExams(ctx.Request.Context(), ports.ListAllExamsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		examCount, err := h.App.CountExams(ctx)
		if err != nil {
			return
		}

		props := components.ExamTableProps{
			Title:      "All Exams",
			Exams:      exams,
			TotalCount: examCount,
			Limit:      int32(limit),
			Offset:     int32(offset),
			BaseURL:    "/admin/dashboard/exams",
		}

		// Check if it's an HTMX request
		if ctx.GetHeader("HX-Request") != "" {
			render.Render(ctx, components.ExamTableContainer(props))
			return
		}

		username, fullname := parseUsername(ctx)

		params := page.AllExamsPageParam{
			FullName:   fullname,
			Username:   username,
			Exams:      exams,
			TotalExams: examCount,
			Limit:      int32(limit),
			Offset:     int32(offset),
		}
		render.Render(ctx, pages.BasePage("All Exams", page.AllExamsPage(params)))
	}
}

func (h *UserHandler) SearchExams() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		search := ctx.PostForm("search")
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
		if helpers.IsTrimmedEmpty(search) {
			exams, err := h.App.ListAllExams(ctx, ports.ListAllExamsParams{
				Limit:  int32(limit),
				Offset: int32(offset),
			})
			if err != nil {
				exams = []domain.Exam{}
			}
			totalCount, _ := h.App.CountExams(ctx)
			ctx.Header("Content-Type", "text/html")
			render.Render(ctx, components.ExamTableContainer(components.ExamTableProps{
				Exams:      exams,
				TotalCount: totalCount,
				Limit:      int32(limit),
				Offset:     int32(offset),
				BaseURL:    "/admin/dashboard/exams",
			}))
			return
		}

		exams, err := h.App.SearchExams(ctx, ports.SearchExamsParams{
			Search: search,
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			render.Render(ctx, components.ExamTableRows([]domain.Exam{}))
			return
		}

		if exams == nil {
			exams = []domain.Exam{}
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.ExamTableContainer(components.ExamTableProps{
			Exams:      exams,
			TotalCount: 0, // Search results often don't show full pagination
			Limit:      int32(limit),
			Offset:     int32(offset),
			BaseURL:    "/admin/dashboard/exams",
		}))
	}
}

func (h *UserHandler) SoftDeleteExam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
			return
		}

		if err := h.App.SoftDeleteExam(ctx, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam"})
			return
		}

		// After deletion, we re-render the table container.
		// For simplicity, we just fetch the first page.
		exams, err := h.App.ListAllExams(ctx, ports.ListAllExamsParams{Limit: 10, Offset: 0})
		if err != nil {
			exams = []domain.Exam{}
		}
		totalCount, _ := h.App.CountExams(ctx)

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.ExamTableContainer(components.ExamTableProps{
			Exams:      exams,
			TotalCount: totalCount,
			Limit:      10,
			Offset:     0,
			BaseURL:    "/admin/dashboard/exams",
		}))
	}
}

func (h *UserHandler) AllSubjectsPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjects, err := h.App.ListAllSubjects(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		username, fullname := parseUsername(ctx)
		render.Render(ctx, pages.BasePage("All Subjects", page.AllSubjectsPage(page.AllSubjectsPageParam{
			Subjects: subjects,
			Username: username,
			FullName: fullname,
		})))
	}
}

func (h *UserHandler) SearchSubjects() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		search := ctx.PostForm("search")
		subjects, err := h.App.SearchSubjects(ctx.Request.Context(), search)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		render.Render(ctx, components.SubjectTable(subjects))
	}
}

func (h *UserHandler) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
	}
}

func (h *UserHandler) EditSubjectPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), uuid.MustParse(subjectID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		render.Render(ctx, pages.BasePage("Edit Subject", page.EditSubjectPage(page.EditSubjectPageParam{
			Subject: subject,
		})))
	}
}
