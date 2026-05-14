package handler

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/components/toast"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/components"
	"uot-exam/web/templates/pages/admin_dashboard/page"
	"uot-exam/web/templates/render"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

var startTime = time.Now()

func (h *UserHandler) AdminDashboardMainRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.App.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{Limit: 6, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		username, fullname, userID := parseUsername(ctx)
		exams, err := h.App.ListAllExams(ctx.Request.Context(), ports.ListAllExamsParams{Limit: 6, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userCount, err := h.App.CountUsers(ctx)
		if err != nil {
			return
		}

		examCount, err := h.App.CountExams(ctx.Request.Context())
		if err != nil {
			return
		}

		questionCount, err := h.App.CountQuestions(ctx.Request.Context())
		if err != nil {
			return
		}

		params := page.AdminDashboardParam{
			Users:               users,
			Username:            username,
			FullName:            fullname,
			UserID:              userID.String(),
			Exams:               exams,
			TotalUsers:          fmt.Sprintf("%d", userCount),
			TotalUsersCount:     userCount,
			TotalExams:          fmt.Sprintf("%d", examCount),
			TotalExamsCount:     examCount,
			TotalQuestions:      fmt.Sprintf("%d", questionCount),
			TotalQuestionsCount: questionCount,
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AdminMainPage(params)))
	}
}

func (h *UserHandler) LandingPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var systemName, adminEmail string
		_ = h.App.DB().Pool.QueryRow(ctx.Request.Context(), "SELECT value FROM system_settings WHERE key = 'system_name'").Scan(&systemName)
		_ = h.App.DB().Pool.QueryRow(ctx.Request.Context(), "SELECT value FROM system_settings WHERE key = 'admin_email'").Scan(&adminEmail)

		if systemName == "" {
			systemName = "UoT Examination System"
		}
		if adminEmail == "" {
			adminEmail = "admin@uotechnology.edu.iq"
		}

		render.Render(ctx, pages.BasePage(systemName, page.LandingPage(systemName, adminEmail)))
	}
}

func (h *UserHandler) AdminLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		count, err := h.App.CountAdmins(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if count == 0 {
			render.Render(ctx, pages.BasePage("Initial Setup", page.InitialSetupPage()))
			return
		}

		render.Render(ctx, pages.BasePage("Admin Login", page.LoginPage("")))
	}
}

func (h *UserHandler) UserPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitStr := ctx.DefaultQuery("limit", "12")
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

		deletedUsers, err := h.App.ListDeletedUsers(ctx.Request.Context(), ports.ListDeletedUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		username, fullname, userID := parseUsername(ctx)
		userCount, err = h.App.CountUsers(ctx)
		if err != nil {
			return
		}

		deletedUsers, _ = h.App.ListDeletedUsers(ctx.Request.Context(), ports.ListDeletedUsersParams{Limit: 12, Offset: 0})
		deletedCount, _ := h.App.CountDeletedUsers(ctx)
		instructorCount, _ := h.App.CountInstructors(ctx.Request.Context())
		studentCount, _ := h.App.CountStudents(ctx.Request.Context())

		params := page.AdminDashboardParam{
			Users:             users,
			Username:          username,
			FullName:          fullname,
			UserID:            userID.String(),
			TotalUsers:        fmt.Sprintf("%d", userCount),
			TotalUsersCount:   userCount,
			TotalInstructors:  fmt.Sprintf("%d", instructorCount),
			TotalStudents:     fmt.Sprintf("%d", studentCount),
			DeletedUsers:      deletedUsers,
			DeletedUsersCount: deletedCount,
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AllUsers(params)))
	}
}

func parseUsername(ctx *gin.Context) (string, string, uuid.UUID) {
	usernameVal, exists := ctx.Get("username")
	if !exists {
		return "", "", uuid.Nil
	}
	username := usernameVal.(string)

	fullnameVal, exists := ctx.Get("fullname")
	if !exists {
		return "", "", uuid.Nil
	}
	fullname := fullnameVal.(string)

	userIDVal, exists := ctx.Get("userID")
	if !exists {
		return "", "", uuid.Nil
	}
	userID := uuid.Must(uuid.Parse(userIDVal.(string)))

	return username, fullname, userID
}

func (h *UserHandler) AllExamsPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitStr := ctx.DefaultQuery("limit", "12")
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

		username, fullname, userID := parseUsername(ctx)

		params := page.AllExamsPageParam{
			FullName:   fullname,
			Username:   username,
			UserID:     userID.String(),
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
			render.Render(ctx, components.ExamTableGrid([]domain.Exam{}))
			return
		}

		if exams == nil {
			exams = []domain.Exam{}
		}

		totalCount, err := h.App.CountSearchExams(ctx, search)
		if err != nil {
			totalCount = 0
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.ExamTableContainer(components.ExamTableProps{
			Exams:      exams,
			TotalCount: totalCount,
			Limit:      int32(limit),
			Offset:     int32(offset),
			BaseURL:    "/admin/exams/search?search=" + url.QueryEscape(search),
			Search:     true,
			SearchAPI:  "/admin/exams/search",
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

		exam, err := h.App.GetExamByID(ctx.Request.Context(), id)
		if err == nil && exam.Status != domain.ExamStatusDraft {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Only draft exams can be deleted"})
			return
		}

		if err := h.App.SoftDeleteExam(ctx, id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam"})
			return
		}

		if ctx.GetHeader("HX-Target") == "toast-container" {
			ctx.Header("HX-Redirect", "/admin/dashboard/exams")
			return
		}

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
		limit := ctx.DefaultQuery("limit", "12")
		limitInt, _ := strconv.Atoi(limit)
		offset := ctx.DefaultQuery("offset", "0")
		offsetInt, _ := strconv.Atoi(offset)

		subjects, err := h.App.ListAllSubjects(ctx.Request.Context(), int32(limitInt), int32(offsetInt))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalCount, err := h.App.CountSubjects(ctx.Request.Context())
		if err != nil {
			totalCount = 0
		}

		username, fullname, userID := parseUsername(ctx)

		if ctx.GetHeader("HX-Request") != "" {
			render.Render(ctx, components.SubjectTableContainer(components.SubjectTableContainerProps{
				Subjects:   subjects,
				TotalCount: totalCount,
				Limit:      int32(limitInt),
				Offset:     int32(offsetInt),
				BaseURL:    "/admin/dashboard/subjects",
			}))
			return
		}

		render.Render(ctx, pages.BasePage("All Subjects", page.AllSubjectsPage(page.AllSubjectsPageParam{
			Subjects:   subjects,
			Username:   username,
			FullName:   fullname,
			UserID:     userID.String(),
			TotalCount: totalCount,
			Limit:      int32(limitInt),
			Offset:     int32(offsetInt),
		})))
	}
}

func (h *UserHandler) SearchSubjects() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit := ctx.DefaultQuery("limit", "12")
		limitInt, _ := strconv.Atoi(limit)
		offset := ctx.DefaultQuery("offset", "0")
		offsetInt, _ := strconv.Atoi(offset)

		search := strings.TrimSpace(ctx.PostForm("search"))
		if search == "" {
			search = strings.TrimSpace(ctx.Query("search"))
		}

		var (
			subjects   []domain.Subject
			totalCount int64
			baseURL    string
		)

		if search == "" {
			// Empty search behaves like the normal list.
			baseURL = "/admin/dashboard/subjects"
			var err error
			subjects, err = h.App.ListAllSubjects(ctx.Request.Context(), int32(limitInt), int32(offsetInt))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount, err = h.App.CountSubjects(ctx.Request.Context())
			if err != nil {
				totalCount = 0
			}
		} else {
			baseURL = "/admin/dashboard/subjects/search?search=" + url.QueryEscape(search)

			var err error
			subjects, err = h.App.SearchSubjects(ctx.Request.Context(), search, int32(limitInt), int32(offsetInt))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount, err = h.App.CountSearchSubjects(ctx.Request.Context(), search)
			if err != nil {
				totalCount = 0
			}
		}

		render.Render(ctx, components.SubjectTableContainer(components.SubjectTableContainerProps{
			Subjects:   subjects,
			TotalCount: totalCount,
			Limit:      int32(limitInt),
			Offset:     int32(offsetInt),
			BaseURL:    baseURL,
		}))
	}
}

func (h *UserHandler) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func (h *UserHandler) EditSubjectPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		username, fullname, userID := parseUsername(ctx)

		parsedID, err := uuid.Parse(subjectID)
		if err != nil {
			ctx.HTML(http.StatusNotFound, "404.html", nil)
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), parsedID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Subject not found", toast.VariantError)
			return
		}

		instructors, err := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), parsedID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to fetch instructors", toast.VariantError)
			return
		}

		exams, err := h.App.ListExamsBySubject(ctx.Request.Context(), parsedID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to fetch exams", toast.VariantError)
			return
		}
		allInstructors, err := h.App.ListAllInstructors(ctx.Request.Context(), ports.ListAllInstructorsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to fetch instructors", toast.VariantError)
			return
		}
		availableInstructors := filterAvailableUsers(allInstructors, instructors)

		limitStr := ctx.DefaultQuery("limit", "12")
		offsetStr := ctx.DefaultQuery("offset", "0")
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		students, err := h.App.ListStudentsBySubjectIDPaginated(ctx.Request.Context(), parsedID, int32(limit), int32(offset))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to fetch students", toast.VariantError)
			return
		}
		studentCount, err := h.App.CountStudentsBySubjectID(ctx.Request.Context(), parsedID)
		if err != nil {
			studentCount = 0
		}
		allStudents, err := h.App.ListAllStudents(ctx.Request.Context(), ports.ListAllStudentsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to fetch students", toast.VariantError)
			return
		}

		allEnrolledStudents, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), parsedID)
		availableStudents := filterAvailableUsers(allStudents, allEnrolledStudents)

		analytics, err := h.App.GetSubjectAnalytics(ctx.Request.Context(), parsedID)
		if err != nil {
			analytics = domain.SubjectAnalytics{}
		}

		attempts, err := h.App.ListOverallAttemptsBySubject(ctx.Request.Context(), parsedID)
		if err != nil {
			attempts = []domain.SubjectAttempt{}
		}

		render.Render(ctx, pages.BasePage("Edit Subject", page.EditSubjectPage(page.EditSubjectPageParam{
			Subject:         subject,
			Username:        username,
			FullName:        fullname,
			Instructors:     instructors,
			Exams:           exams,
			UserID:          userID,
			AllInstructors:  availableInstructors,
			Students:        students,
			AllStudents:     availableStudents,
			StudentCount:    studentCount,
			Analytics:       analytics,
			OverallAttempts: attempts,
			Limit:           int32(limit),
			Offset:          int32(offset),
		})))
	}
}

func (h *UserHandler) SearchExamsBySubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
			return
		}

		search := strings.TrimSpace(ctx.PostForm("search"))
		if search == "" {
			search = strings.TrimSpace(ctx.Query("search"))
		}

		limitStr := ctx.DefaultQuery("limit", "10")
		limit, _ := strconv.Atoi(limitStr)
		offsetStr := ctx.DefaultQuery("offset", "0")
		offset, _ := strconv.Atoi(offsetStr)

		var exams []domain.Exam

		if search == "" {
			exams, err = h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
			if err != nil {
				exams = []domain.Exam{}
			}
		} else {
			exams, err = h.App.SearchExamsBySubject(ctx.Request.Context(), ports.SearchExamsBySubjectParams{
				SubjectID: subjectID,
				Search:    search,
				Limit:     int32(limit),
				Offset:    int32(offset),
			})
			if err != nil {
				exams = []domain.Exam{}
			}
		}

		ctx.Header("Content-Type", "text/html")
		// Determine baseURL based on user role
		baseURL := "/admin/dashboard/exam/"
		if role, exists := ctx.Get("role"); exists {
			if domain.UserRole(role.(string)) == domain.RoleInstructor {
				baseURL = "/instructor/exam/"
			}
		}

		render.Render(ctx, components.ExamTableGrid(exams, baseURL))
	}
}

func (h *UserHandler) CreateExam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		subjectIDStr := ctx.PostForm("subject_id")
		// Get userID from context (set by AuthenticationMiddleware)
		userIDVal, exists := ctx.Get("userID")
		if !exists {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Session expired. Please log out and log in again.", toast.VariantError)
			return
		}

		createdBy, err := uuid.Parse(userIDVal.(string))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid session user ID", toast.VariantError)
			return
		}
		description := ctx.PostForm("description")
		totalMarksStr := ctx.PostForm("total_marks")
		passScoreStr := ctx.PostForm("pass_score")

		if helpers.IsTrimmedEmpty(title) || helpers.IsTrimmedEmpty(subjectIDStr) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Missing required fields", toast.VariantError)
			return
		}

		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid subject ID", toast.VariantError)
			return
		}

		totalMarks, err := strconv.ParseFloat(totalMarksStr, 64)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid total marks", toast.VariantError)
			return
		}

		passScore, _ := strconv.ParseFloat(passScoreStr, 64)
		if passScoreStr == "" {
			passScore = totalMarks / 2.0
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		if err == nil && subject.Status == domain.SubjectStatusPublished {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Subject is published and cannot be modified", toast.VariantError)
			return
		}

		_, err = h.App.CreateExam(ctx, ports.CreateExamParams{
			Title:       title,
			SubjectID:   subjectID,
			CreatedBy:   createdBy,
			Description: &description,
			TotalMarks:  totalMarks,
			PassScore:   passScore,
			Status:      domain.ExamStatusDraft,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Failed to create exam : "+err.Error(), toast.VariantError)
			return
		}

		// Return updated exam grid
		exams, err := h.App.ListExamsBySubject(ctx, subjectID)
		if err != nil {
			exams = []domain.Exam{}
		}

		render.Render(ctx, components.ExamTableGrid(exams))
		helpers.Toast(ctx, "Create Exam Success", "Exam created successfully", toast.VariantSuccess)
	}
}

func (h *UserHandler) EditSubjectInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		name := ctx.PostForm("subject_name")
		description := ctx.PostForm("subject_description")

		if helpers.IsTrimmedEmpty(name) {
			helpers.Toast(ctx, "Edit Subject Failed", "Subject name cannot be empty", toast.VariantError)
			return
		}

		if len(description) > 255 {
			helpers.Toast(ctx, "Edit Subject Failed", "Subject description cannot be longer than 255 characters", toast.VariantError)
			return
		}

		durationStr := ctx.PostForm("duration_minutes")
		passScoreStr := ctx.PostForm("pass_score")
		statusStr := ctx.PostForm("status")

		passScore, _ := strconv.ParseFloat(passScoreStr, 64)

		var duration int32
		if durationStr != "" {
			if strings.Contains(durationStr, ":") {
				parts := strings.Split(durationStr, ":")
				if len(parts) == 2 {
					hour, _ := strconv.Atoi(parts[0])
					minute, _ := strconv.Atoi(parts[1])
					duration = int32(hour)*60 + int32(minute)
				}
			} else {
				d, _ := strconv.Atoi(durationStr)
				duration = int32(d)
			}
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), uuid.MustParse(id))
		if err != nil {
			helpers.Toast(ctx, "Edit Subject Failed", "Subject not found", toast.VariantError)
			return
		}

		if subject.Status == domain.SubjectStatusPublished {
			helpers.Toast(ctx, "Edit Subject Failed", "Subject is published and cannot be modified", toast.VariantError)
			return
		}

		if passScore == 0 && subject.TotalMarks > 0 {
			passScore = subject.TotalMarks / 2.0
		}

		_, err = h.App.UpdateSubject(ctx, domain.Subject{
			ID:              uuid.MustParse(id),
			Title:           name,
			Description:     &description,
			DurationMinutes: duration,
			PassScore:       passScore,
			Status:          domain.SubjectStatus(statusStr),
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to update subject : "+err.Error(), toast.VariantError)
			return
		}

		helpers.Toast(ctx, "Edit Subject Success", "Subject updated successfully", toast.VariantSuccess)
	}
}

func (h *UserHandler) PublishSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")

		subjectID, err := uuid.Parse(idStr)
		if err != nil {
			helpers.Toast(ctx, "Publish Failed", "Invalid subject ID format", toast.VariantError)
			return
		}

		// check if it exists
		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		if err != nil {
			helpers.Toast(ctx, "Publish Failed", "Subject not found", toast.VariantError)
			return
		}

		if subject.Status == domain.SubjectStatusPublished {
			helpers.Toast(ctx, "Notice", "Subject is already published", toast.VariantWarning)
			return
		}

		if subject.PassScore <= 0 {
			helpers.Toast(ctx, "Publish Failed", "Pass score cannot be zero. Please set a pass score before publishing.", toast.VariantError)
			return
		}

		// check for students
		students, err := h.App.ListStudentsBySubjectID(ctx.Request.Context(), subjectID)
		if err != nil || len(students) == 0 {
			helpers.Toast(ctx, "Publish Failed", "Subject must have at least one student assigned.", toast.VariantError)
			return
		}

		// check for exams
		exams, err := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
		if err != nil || len(exams) == 0 {
			helpers.Toast(ctx, "Publish Failed", "Subject must have at least one exam.", toast.VariantError)
			return
		}

		// update the subject
		err = h.App.PublishSubject(ctx.Request.Context(), subjectID)
		if err != nil {
			helpers.Toast(ctx, "Publish Failed", "Failed to update subject status: "+err.Error(), toast.VariantError)
			return
		}

		// update draft exams
		err = h.App.PublishDraftExamsBySubject(ctx.Request.Context(), subjectID)
		if err != nil {
			helpers.Toast(ctx, "Publish Failed", "Failed to publish exams: "+err.Error(), toast.VariantError)
			return
		}

		// HTMX Redirect to refresh lockdown cleanly
		ctx.Header("HX-Redirect", "/admin/dashboard/subject/"+idStr)
		helpers.Toast(ctx, "Publish Success", "Subject and its draft exams have been published", toast.VariantSuccess)
	}
}

func (h *UserHandler) CreateSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.PostForm("title")
		description := ctx.PostForm("description")

		if helpers.IsTrimmedEmpty(name) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Subject Failed", "Subject name cannot be empty", toast.VariantError)
			return
		}

		if len(description) > 255 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Subject Failed", "Subject description cannot be longer than 255 characters", toast.VariantError)
			return
		}

		durationStr := ctx.PostForm("duration_minutes")
		statusStr := ctx.PostForm("status")
		passScore := 0.0

		var duration int32
		if durationStr != "" {
			if strings.Contains(durationStr, ":") {
				parts := strings.Split(durationStr, ":")
				if len(parts) == 2 {
					hour, _ := strconv.Atoi(parts[0])
					minute, _ := strconv.Atoi(parts[1])
					duration = int32(hour)*60 + int32(minute)
				}
			} else {
				d, _ := strconv.Atoi(durationStr)
				duration = int32(d)
			}
		}

		_, err := h.App.CreateSubject(ctx, domain.Subject{
			Title:           name,
			Description:     &description,
			DurationMinutes: duration,
			PassScore:       passScore,
			Status:          domain.SubjectStatus(statusStr),
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Subject Failed", "Failed to create subject : "+err.Error(), toast.VariantError)
			return
		}
		limit := ctx.DefaultQuery("limit", "12")
		limitInt, _ := strconv.Atoi(limit)
		offset := ctx.DefaultQuery("offset", "0")
		offsetInt, _ := strconv.Atoi(offset)
		subjects, err := h.App.ListAllSubjects(ctx.Request.Context(), int32(limitInt), int32(offsetInt))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Subject Failed", "Failed to list subjects", toast.VariantError)
			return
		}
		totalCount, err := h.App.CountSubjects(ctx.Request.Context())
		if err != nil {
			totalCount = 0
		}

		render.Render(ctx, components.SubjectTableContainer(components.SubjectTableContainerProps{
			Subjects:   subjects,
			TotalCount: totalCount,
			Limit:      int32(limitInt),
			Offset:     int32(offsetInt),
			BaseURL:    "/admin/dashboard/subjects",
		}))
		helpers.Toast(ctx, "Create Subject Success", "Subject created successfully", toast.VariantSuccess)
	}
}

func (h *UserHandler) DeleteSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if helpers.IsTrimmedEmpty(id) {
			helpers.Toast(ctx, "Delete Subject Failed", "Subject ID cannot be empty", toast.VariantError)
			return
		}
		subjectID, err := uuid.Parse(id)
		if err != nil {
			helpers.Toast(ctx, "Delete Subject Failed", "Invalid subject ID", toast.VariantError)
			return
		}
		_, err = h.App.DeleteSubject(ctx, subjectID)
		if err != nil {
			helpers.Toast(ctx, "Delete Subject Failed", "Failed to delete subject", toast.VariantError)
			return
		}
		ctx.Header("HX-Redirect", "/admin/dashboard/subjects")
	}
}

func (h *UserHandler) EditExamPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		examID := ctx.Param("id")
		if helpers.IsTrimmedEmpty(examID) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Exam ID cannot be empty", toast.VariantError)
			return
		}
		examIDUUID, err := uuid.Parse(examID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Invalid exam ID", toast.VariantError)
			return
		}
		exam, err := h.App.GetExamByID(ctx.Request.Context(), examIDUUID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Failed to get exam", toast.VariantError)
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), exam.SubjectID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Failed to get subject", toast.VariantError)
			return
		}

		questions, err := h.App.ListQuestionsByExam(ctx.Request.Context(), examIDUUID)
		if err != nil {
			questions = []domain.Question{}
		}

		attempts, err := h.App.ListAttemptsByExam(ctx.Request.Context(), examIDUUID)
		if err != nil {
			attempts = []domain.ExamAttempt{}
		}

		analytics, err := h.App.GetExamAnalytics(ctx.Request.Context(), examIDUUID)
		if err != nil {
			analytics = domain.ExamAnalytics{}
		}

		var choices []domain.Choice
		for i := range questions {
			choices, err = h.App.ListChoicesByQuestion(ctx.Request.Context(), questions[i].ID)
			if err != nil {
				choices = []domain.Choice{}
			}
			questions[i].Choices = choices
		}

		username, fullname, userID := parseUsername(ctx)
		render.Render(ctx, pages.BasePage("Edit Exam", page.EditExamPage(page.EditExamPageParam{
			Username:  username,
			FullName:  fullname,
			UserID:    userID.String(),
			Exam:      exam,
			Subject:   subject,
			Questions: questions,
			Attempts:  attempts,
			Analytics: analytics,
		})))
	}
}

func (h *UserHandler) EditExamInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		name := ctx.PostForm("exam_name")
		description := ctx.PostForm("exam_description")
		totalMarksStr := ctx.PostForm("exam_total_marks")
		passScoreStr := ctx.PostForm("exam_pass_score")
		statusStr := ctx.PostForm("exam_status")

		if helpers.IsTrimmedEmpty(idStr) || helpers.IsTrimmedEmpty(name) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Exam ID and Name are required", toast.VariantError)
			return
		}

		examID, err := uuid.Parse(idStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Invalid exam ID", toast.VariantError)
			return
		}

		exam, err := h.App.GetExamByID(ctx.Request.Context(), examID)
		if err == nil && exam.Status != domain.ExamStatusDraft {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Only draft exams can be modified", toast.VariantError)
			return
		}

		totalMarks, _ := strconv.ParseFloat(totalMarksStr, 64)
		passScore, _ := strconv.ParseFloat(passScoreStr, 64)

		if passScore == exam.PassScore && totalMarks != exam.TotalMarks {
			passScore = totalMarks / 2.0
		}

		_, err = h.App.UpdateExam(ctx, ports.UpdateExamParams{
			ID:          examID,
			Title:       name,
			Description: &description,
			TotalMarks:  totalMarks,
			PassScore:   passScore,
			Status:      domain.ExamStatus(statusStr),
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Failed to update exam: "+err.Error(), toast.VariantError)
			return
		}

		helpers.Toast(ctx, "Edit Exam Success", "Exam updated successfully", toast.VariantSuccess)
	}
}

func (h *UserHandler) PreviewQuestionText() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		questionText := ctx.PostForm("question_text")
		questionTitle := ctx.PostForm("question_title")
		if helpers.IsTrimmedEmpty(questionTitle) {
			questionTitle = `\[\text{Question Title}\]`
		}
		if helpers.IsTrimmedEmpty(questionText) {
			questionText = `\[\text{Question Text}\]`
		}
		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.QuestionPreview("question-preview", questionTitle, questionText))
	}
}

func (h *UserHandler) GetQuestionForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		questionType := ctx.PostForm("question_type")
		ctx.Header("Content-Type", "text/html")
		switch domain.QuestionType(questionType) {
		case domain.QuestionTypeText:
			render.Render(ctx, components.TextQuestionForm(domain.Question{}, "question-preview"))
		case domain.QuestionTypeCode:
			render.Render(ctx, components.CodeQuestionForm(domain.Question{}, "question-preview"))
		case domain.QuestionTypeImage:
			render.Render(ctx, components.ImageQuestionForm(domain.Question{}, true, "question-preview"))
		default:
			render.Render(ctx, components.TextQuestionForm(domain.Question{}, "question-preview"))
		}
	}
}

func (h *UserHandler) PreviewQuestionChoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		triggerName := ctx.GetHeader("HX-Trigger-Name")

		if triggerName == "" {
			ctx.Status(http.StatusBadRequest)
			return
		}

		inputValue := ctx.PostForm(triggerName)

		if helpers.IsTrimmedEmpty(inputValue) {
			displayTitle := strings.Replace(triggerName, "choice_", "Choice ", 1)
			inputValue = fmt.Sprintf(`$\text{%s}$`, displayTitle)
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.QuestionChoicesPreview(inputValue))
	}
}

func (h *UserHandler) CreateQuestion() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		questionTitle := ctx.PostForm("question_title")
		questionText := ctx.PostForm("question_text")
		questionType := ctx.PostForm("question_type")
		correctChoice := ctx.PostForm("correct_choice")

		var questionImageURL string
		if questionType == string(domain.QuestionTypeImage) {
			questionImage, err := ctx.FormFile("question_image")
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Create Question Failed", "Failed to get question image", toast.VariantError)
				return
			}

			helpers.ValidateImage(ctx, questionImage)

			hashedFileName := helpers.HashFileName(questionImage.Filename)
			questionImageURL, err = h.App.UploadQuestionImage(questionImage, hashedFileName)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Create Question Failed", "Failed to upload question image", toast.VariantError)
				return
			}
		}

		choices := []domain.Choice{}

		for i := 1; i <= 4; i++ {
			choice := ctx.PostForm("choice_" + strconv.Itoa(i))
			if !helpers.IsTrimmedEmpty(choice) {
				choices = append(choices, domain.Choice{
					ChoiceText: choice,
					IsCorrect:  fmt.Sprintf("choice_%d", i) == correctChoice,
				})
			}
		}

		examID := ctx.Param("id")
		if helpers.IsTrimmedEmpty(examID) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "Exam ID cannot be empty", toast.VariantError)
			return
		}

		if helpers.IsTrimmedEmpty(questionTitle) || (questionType != string(domain.QuestionTypeImage) && helpers.IsTrimmedEmpty(questionText)) ||
			helpers.IsTrimmedEmpty(questionType) || len(choices) == 0 ||
			helpers.IsTrimmedEmpty(correctChoice) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "All fields are required", toast.VariantError)
			return
		}

		questionMarksFloat := 1.0

		parsedUUID, err := uuid.Parse(examID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "Invalid exam ID", toast.VariantError)
			return
		}

		questionChecksum, err := helpers.BuildQuestionChecksum(examID, domain.Question{
			QuestionTitle: questionTitle,
			QuestionText:  questionText,
			QuestionType:  questionType,
			Marks:         questionMarksFloat,
			Choices:       choices,
		})

		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", err.Error(), toast.VariantError)
			return
		}

		exists, err := h.App.GetQuestionByChecksum(ctx, questionChecksum)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "Failed to get question: "+err.Error(), toast.VariantError)
			return
		}

		if exists {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "Question already exists", toast.VariantError)
			return
		}

		question, err := h.App.CreateQuestion(ctx, ports.CreateQuestionParams{
			ExamID:        parsedUUID,
			QuestionTitle: questionTitle,
			QuestionText:  questionText,
			QuestionType:  domain.QuestionType(questionType),
			Marks:         questionMarksFloat,
			ImageURL:      questionImageURL,
			Checksum:      questionChecksum,
		})

		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "Failed to create question: "+err.Error(), toast.VariantError)
			return
		}

		for _, choice := range choices {
			_, err = h.App.CreateChoice(ctx, ports.CreateChoiceParams{
				QuestionID: question.ID,
				ChoiceText: choice.ChoiceText,
				IsCorrect:  choice.IsCorrect,
			})
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Create Question Failed", "Failed to create choice: "+err.Error(), toast.VariantError)
				return
			}
		}

		// Reload question with choices
		question.Choices, err = h.App.ListChoicesByQuestion(ctx.Request.Context(), question.ID)
		if err != nil {
			question.Choices = []domain.Choice{}
		}

		helpers.Toast(ctx, "Create Question Success", "Question created successfully", toast.VariantSuccess)
		render.Render(ctx, components.QuestionList(components.QuestionListProps{
			Questions: []domain.Question{
				question,
			},
		}))
	}
}

func (h *UserHandler) UploadQuestionCSV() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		examID := ctx.Param("id")

		if helpers.IsTrimmedEmpty(examID) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Upload Question CSV Failed", "Exam ID cannot be empty", toast.VariantError)
			return
		}

		parsedUUID, err := uuid.Parse(examID)

		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Upload Question CSV Failed", "Invalid exam ID", toast.VariantError)
			return
		}

		file, err := helpers.ParseCSVFile(ctx, "file")
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to parse CSV file", toast.VariantError)
			return
		}

		questions, err := helpers.MapQuestionCSVToStruct(file)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to parse CSV file", toast.VariantError)
			return
		}

		for _, question := range questions {

			questionChecksum, err := helpers.BuildQuestionChecksum(examID, question)

			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to generate question checksum", toast.VariantError)
				return
			}

			existingQuestion, err := h.App.GetQuestionByChecksum(ctx, questionChecksum)

			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to get question: "+err.Error(), toast.VariantError)
				return
			}

			if existingQuestion {
				continue
			}

			createdQuestion, err := h.App.CreateQuestion(ctx, ports.CreateQuestionParams{
				ExamID:        parsedUUID,
				QuestionTitle: question.QuestionTitle,
				QuestionText:  question.QuestionText,
				QuestionType:  domain.QuestionType(question.QuestionType),
				Marks:         question.Marks,
				ImageURL:      "",
				Checksum:      questionChecksum,
			})

			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to create question: "+err.Error(), toast.VariantError)
				return
			}
			for _, choice := range question.Choices {
				_, err = h.App.CreateChoice(ctx, ports.CreateChoiceParams{
					QuestionID: createdQuestion.ID,
					ChoiceText: choice.ChoiceText,
					IsCorrect:  choice.IsCorrect,
				})
				if err != nil {
					ctx.Header("HX-Reswap", "none")
					helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to create choice: "+err.Error(), toast.VariantError)
					return
				}
			}
		}

		questions, err = h.App.ListQuestionsByExam(ctx, parsedUUID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to reload questions: "+err.Error(), toast.VariantError)
			return
		}

		for i := range questions {
			questions[i].Choices, err = h.App.ListChoicesByQuestion(ctx, questions[i].ID)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Upload Question CSV Failed", "Failed to reload questions: "+err.Error(), toast.VariantError)
				return
			}
		}

		helpers.Toast(ctx, "Upload Question CSV Success", "Question uploaded successfully", toast.VariantSuccess)
		render.Render(ctx, components.QuestionList(components.QuestionListProps{
			Questions: questions,
		}))
	}
}

func (h *UserHandler) UploadQuestionCSVRandom() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		examID := ctx.Param("id")

		if helpers.IsTrimmedEmpty(examID) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Exam ID cannot be empty", toast.VariantError)
			return
		}

		parsedUUID, err := uuid.Parse(examID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Invalid exam ID", toast.VariantError)
			return
		}

		// Parse the requested count
		countStr := ctx.PostForm("count")
		if helpers.IsTrimmedEmpty(countStr) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Number of questions is required", toast.VariantError)
			return
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count < 1 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Number of questions must be a positive integer", toast.VariantError)
			return
		}

		// Parse and validate CSV file
		file, err := helpers.ParseCSVFile(ctx, "file")
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to parse CSV file: "+err.Error(), toast.VariantError)
			return
		}

		allQuestions, err := helpers.MapQuestionCSVToStruct(file)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to parse CSV file: "+err.Error(), toast.VariantError)
			return
		}

		if len(allQuestions) == 0 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "No valid questions found in the CSV file", toast.VariantError)
			return
		}

		// Clamp count to available questions
		if count > len(allQuestions) {
			count = len(allQuestions)
		}

		// Randomly shuffle and pick the first 'count' questions
		rand.Shuffle(len(allQuestions), func(i, j int) {
			allQuestions[i], allQuestions[j] = allQuestions[j], allQuestions[i]
		})
		selected := allQuestions[:count]

		addedCount := 0
		for _, question := range selected {

			questionChecksum, err := helpers.BuildQuestionChecksum(examID, question)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to generate question checksum", toast.VariantError)
				return
			}

			existingQuestion, err := h.App.GetQuestionByChecksum(ctx, questionChecksum)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to check question: "+err.Error(), toast.VariantError)
				return
			}

			if existingQuestion {
				continue
			}

			createdQuestion, err := h.App.CreateQuestion(ctx, ports.CreateQuestionParams{
				ExamID:        parsedUUID,
				QuestionTitle: question.QuestionTitle,
				QuestionText:  question.QuestionText,
				QuestionType:  domain.QuestionType(question.QuestionType),
				Marks:         question.Marks,
				ImageURL:      "",
				Checksum:      questionChecksum,
			})

			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to create question: "+err.Error(), toast.VariantError)
				return
			}

			for _, choice := range question.Choices {
				_, err = h.App.CreateChoice(ctx, ports.CreateChoiceParams{
					QuestionID: createdQuestion.ID,
					ChoiceText: choice.ChoiceText,
					IsCorrect:  choice.IsCorrect,
				})
				if err != nil {
					ctx.Header("HX-Reswap", "none")
					helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to create choice: "+err.Error(), toast.VariantError)
					return
				}
			}
			addedCount++
		}

		// Reload all questions for the exam
		questions, err := h.App.ListQuestionsByExam(ctx, parsedUUID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to reload questions: "+err.Error(), toast.VariantError)
			return
		}

		for i := range questions {
			questions[i].Choices, err = h.App.ListChoicesByQuestion(ctx, questions[i].ID)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Random CSV Upload Failed", "Failed to reload questions: "+err.Error(), toast.VariantError)
				return
			}
		}

		msg := fmt.Sprintf("%d random question(s) added successfully", addedCount)
		helpers.Toast(ctx, "Random CSV Upload Success", msg, toast.VariantSuccess)
		render.Render(ctx, components.QuestionList(components.QuestionListProps{
			Questions: questions,
		}))
	}
}

func (h *UserHandler) DeleteQuestion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		questionID := ctx.Param("question-id")

		if helpers.IsTrimmedEmpty(questionID) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Delete Question Failed", "Question ID cannot be empty", toast.VariantError)
			return
		}

		parsedQuestionUUID, err := uuid.Parse(questionID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Delete Question Failed", "Invalid question ID", toast.VariantError)
			return
		}

		err = h.App.DeleteQuestionAndChoices(ctx, parsedQuestionUUID)

		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Delete Question Failed", "Failed to delete question: "+err.Error(), toast.VariantError)
			return
		}

		examID := ctx.Param("id")
		parsedExamUUID, _ := uuid.Parse(examID)

		questions, err := h.App.ListQuestionsByExam(ctx, parsedExamUUID)
		if err != nil {
			helpers.Toast(ctx, "Delete Question Failed", "Failed to reload questions: "+err.Error(), toast.VariantError)
			return
		}

		for i := range questions {
			questions[i].Choices, err = h.App.ListChoicesByQuestion(ctx, questions[i].ID)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Delete Question Failed", "Failed to reload questions: "+err.Error(), toast.VariantError)
				return
			}
		}

		helpers.Toast(ctx, "Delete Question Success", "Question deleted successfully", toast.VariantSuccess)
		render.Render(ctx, components.QuestionList(components.QuestionListProps{
			Questions: questions,
		}))
	}
}

func (h *UserHandler) ExportExamCSV() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		examID := ctx.Param("id")
		if helpers.IsTrimmedEmpty(examID) {
			return
		}
		parsedExamUUID, err := uuid.Parse(examID)
		if err != nil {
			return
		}

		questions, err := h.App.ListQuestionsByExam(ctx, parsedExamUUID)
		if err != nil {
			questions = []domain.Question{}
		}
		var choices []domain.Choice
		for i := range questions {
			choices, err = h.App.ListChoicesByQuestion(ctx, questions[i].ID)
			if err != nil {
				choices = []domain.Choice{}
			}
			questions[i].Choices = choices
		}

		helpers.ExportExamCSV(ctx, questions)
	}
}
func (h *UserHandler) UpdateQuestion() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Update Question Failed", "Invalid question ID", toast.VariantError)
			return
		}

		questionTitle := ctx.PostForm("question_title")
		questionText := ctx.PostForm("question_text")
		questionType := ctx.PostForm("question_type")
		correctChoice := ctx.PostForm("correct_choice")

		var questionImageURL string
		questionImage, err := ctx.FormFile("question_image")
		if err == nil {
			helpers.ValidateImage(ctx, questionImage)
			hashedFileName := helpers.HashFileName(questionImage.Filename)
			questionImageURL, err = h.App.UploadQuestionImage(questionImage, hashedFileName)
			if err != nil {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Update Question Failed", "Failed to upload question image", toast.VariantError)
				return
			}
		}

		questionMarksFloat := 1.0

		if helpers.IsTrimmedEmpty(questionTitle) || (questionType != string(domain.QuestionTypeImage) && helpers.IsTrimmedEmpty(questionText)) ||
			helpers.IsTrimmedEmpty(questionType) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Update Question Failed", "All fields are required", toast.VariantError)
			return
		}

		existingQuestion, err := h.App.GetQuestionByID(ctx.Request.Context(), id)
		if err != nil {
			helpers.Toast(ctx, "Update Question Failed", "Question not found", toast.VariantError)
			return
		}

		if questionImageURL == "" {
			questionImageURL = existingQuestion.QuestionImage
		}

		choices := []ports.CreateChoiceParams{}

		for i := 1; i <= 4; i++ {
			choiceText := ctx.PostForm("choice_" + strconv.Itoa(i))
			if !helpers.IsTrimmedEmpty(choiceText) {
				choices = append(choices, ports.CreateChoiceParams{
					QuestionID: id,
					ChoiceText: choiceText,
					IsCorrect:  fmt.Sprintf("choice_%d", i) == correctChoice,
				})
			}
		}

		updateParams := ports.UpdateQuestionParams{
			ID:            id,
			QuestionTitle: questionTitle,
			QuestionText:  questionText,
			QuestionType:  domain.QuestionType(questionType),
			Marks:         questionMarksFloat,
			ImageURL:      questionImageURL,
		}

		_, err = h.App.UpdateQuestionWithChoices(ctx.Request.Context(), updateParams, choices)

		if err != nil {
			helpers.Toast(ctx, "Update Question Failed", "Failed to update question : "+err.Error(), toast.VariantError)
			return
		}

		ctx.Header("HX-Refresh", "true")
		helpers.Toast(ctx, "Success", "Question updated successfully", toast.VariantSuccess)
	}
}

func (h *UserHandler) CreateUserCSV() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("csv_file")
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create User CSV Failed", "Failed to upload file", toast.VariantError)
			return
		}
		if file == nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create User CSV Failed", "File is required", toast.VariantError)
			return
		}

		users, err := helpers.ParseUserCSV(file)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create User CSV Failed", "Failed to parse CSV: "+err.Error(), toast.VariantError)
			return
		}

		for _, user := range users {
			_, err = h.App.CreateUser(ctx, ports.CreateUserParams{
				Username: user.Username,
				Password: user.Password,
				Role:     user.Role,
				FullName: user.FullName,
				IsActive: true,
			})
			if err != nil {
				if errors.Is(err, domain.ErrUserAlreadyExists) || strings.Contains(err.Error(), "duplicate key value") {
					continue
				}
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Create User CSV Failed", "Failed to create user: "+err.Error(), toast.VariantError)
				return
			}
		}

		usersList, err := h.App.ListAllUsers(ctx, ports.ListAllUsersParams{
			Limit:  12,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create User CSV Failed", "Failed to reload users: "+err.Error(), toast.VariantError)
			return
		}

		helpers.Toast(ctx, "Success", "User created successfully", toast.VariantSuccess)
		render.Render(ctx, components.UserTableContainer(components.UserTableProps{
			Users:      usersList,
			TotalCount: int64(len(usersList)),
			Limit:      12,
			Offset:     0,
			BaseURL:    "/admin/dashboard/users",
			Search:     true,
			SearchAPI:  "/admin/users/search",
			AddUser:    true,
			AddUserCSV: true,
		}))
	}
}

func filterAvailableUsers(all []domain.User, alreadyAssigned []domain.User) []domain.User {
	if len(alreadyAssigned) == 0 {
		return all
	}
	assigned := make(map[uuid.UUID]struct{}, len(alreadyAssigned))
	for _, u := range alreadyAssigned {
		assigned[u.ID] = struct{}{}
	}
	out := make([]domain.User, 0, len(all))
	for _, u := range all {
		if _, ok := assigned[u.ID]; !ok {
			out = append(out, u)
		}
	}
	return out
}

func parseIDsFromForm(raw string) ([]uuid.UUID, error) {
	parts := strings.Split(raw, ",")
	seen := make(map[uuid.UUID]struct{})
	ids := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := uuid.Parse(p)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func filterUsersBySearch(users []domain.User, q string) []domain.User {
	q = strings.TrimSpace(q)
	if q == "" {
		return users
	}
	needle := strings.ToLower(q)
	out := make([]domain.User, 0, len(users))
	for _, u := range users {
		if strings.Contains(strings.ToLower(u.FullName), needle) || strings.Contains(strings.ToLower(u.Username), needle) {
			out = append(out, u)
		}
	}
	return out
}

func subjectInstructorsTableProps(subjectID uuid.UUID, instructors, allInstructors []domain.User, search string) components.UserTableProps {
	sid := subjectID.String()
	base := "/admin/dashboard/subject/" + sid
	baseURL := base + "/instructors/search"
	if search != "" {
		baseURL += "?search=" + url.QueryEscape(search)
	}
	return components.UserTableProps{
		Users:                      instructors,
		Title:                      "Instructors",
		BaseURL:                    baseURL,
		ID:                         "instructors-table",
		Search:                     true,
		SearchAPI:                  base + "/instructors/search",
		AddUser:                    false,
		AddInstructor:              true,
		Instructors:                allInstructors,
		AddInstructorAPI:           base + "/instructors/assign",
		InstructorAssignSwapTarget: "#instructors-user-table-root",
	}
}

func (h *UserHandler) SearchSubjectInstructors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
			return
		}

		search := strings.TrimSpace(ctx.PostForm("search"))
		if search == "" {
			search = strings.TrimSpace(ctx.Query("search"))
		}

		assigned, err := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), subjectID)
		if err != nil {
			assigned = []domain.User{}
		}
		instructors := filterUsersBySearch(assigned, search)

		allInstructors, err := h.App.ListAllInstructors(ctx.Request.Context(), ports.ListAllInstructorsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			allInstructors = []domain.User{}
		}
		available := filterAvailableUsers(allInstructors, assigned)

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, components.UserTableContainer(subjectInstructorsTableProps(subjectID, instructors, available, search)))
	}
}

func (h *UserHandler) AssignInstructorToSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		rawIDs := ctx.PostForm("instructor_id")
		if subjectID == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Subject ID is required", toast.VariantError)
			return
		}
		if strings.TrimSpace(rawIDs) == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Select at least one instructor", toast.VariantError)
			return
		}
		sid, err := uuid.Parse(subjectID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Invalid subject ID", toast.VariantError)
			return
		}
		instructorIDs, err := parseIDsFromForm(rawIDs)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Invalid instructor selection", toast.VariantError)
			return
		}
		if len(instructorIDs) == 0 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Select at least one instructor", toast.VariantError)
			return
		}
		if err := h.App.AssignInstructorsToSubject(ctx, sid, instructorIDs); err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Failed to assign instructors: "+err.Error(), toast.VariantError)
			return
		}

		instructors, err := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), sid)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Assigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		allInstructors, err := h.App.ListAllInstructors(ctx.Request.Context(), ports.ListAllInstructorsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Instructors Failed", "Assigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		available := filterAvailableUsers(allInstructors, instructors)

		ctx.Header("Content-Type", "text/html")
		msg := "1 instructor assigned successfully"
		if len(instructorIDs) != 1 {
			msg = fmt.Sprintf("%d instructors assigned successfully", len(instructorIDs))
		}
		// Do not set HX-Reswap: none here — that suppresses swapping the refreshed table into hx-target.
		helpers.Toast(ctx, "Success", msg, toast.VariantSuccess)
		props := subjectInstructorsTableProps(sid, instructors, available, "")
		render.Render(ctx, components.UserTableRoot("instructors-user-table-root", props))
	}
}

func (h *UserHandler) UnassignInstructorFromSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		instructorID := ctx.Param("user_id")

		if subjectID == "" || instructorID == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Instructor Failed", "Subject ID and Instructor ID are required", toast.VariantError)
			return
		}

		sid, err := uuid.Parse(subjectID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Instructor Failed", "Invalid subject ID", toast.VariantError)
			return
		}

		iid, err := uuid.Parse(instructorID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Instructor Failed", "Invalid instructor ID", toast.VariantError)
			return
		}

		if err := h.App.UnassignInstructorFromSubject(ctx, sid, iid); err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Instructor Failed", "Failed to unassign instructor: "+err.Error(), toast.VariantError)
			return
		}

		instructors, err := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), sid)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Instructor Failed", "Unassigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}

		allInstructors, err := h.App.ListAllInstructors(ctx.Request.Context(), ports.ListAllInstructorsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Instructor Failed", "Unassigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		available := filterAvailableUsers(allInstructors, instructors)

		ctx.Header("Content-Type", "text/html")
		helpers.Toast(ctx, "Success", "Instructor unassigned successfully", toast.VariantSuccess)
		props := subjectInstructorsTableProps(sid, instructors, available, "")
		render.Render(ctx, components.UserTableRoot("instructors-user-table-root", props))
	}
}

func (h *UserHandler) AssignStudentToSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		rawIDs := ctx.PostForm("student_id")
		if subjectID == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Subject ID is required", toast.VariantError)
			return
		}
		if strings.TrimSpace(rawIDs) == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Select at least one student", toast.VariantError)
			return
		}
		sid, err := uuid.Parse(subjectID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Invalid subject ID", toast.VariantError)
			return
		}
		studentIDs, err := parseIDsFromForm(rawIDs)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Invalid student selection", toast.VariantError)
			return
		}
		if len(studentIDs) == 0 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Select at least one student", toast.VariantError)
			return
		}
		if err := h.App.AssignStudentsToSubject(ctx, sid, studentIDs); err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Failed to assign students: "+err.Error(), toast.VariantError)
			return
		}

		limitStr := ctx.DefaultQuery("limit", "12")
		offsetStr := ctx.DefaultQuery("offset", "0")
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		students, err := h.App.ListStudentsBySubjectIDPaginated(ctx.Request.Context(), sid, int32(limit), int32(offset))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Assigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		totalCount, _ := h.App.CountStudentsBySubjectID(ctx.Request.Context(), sid)
		allStudents, err := h.App.ListAllStudents(ctx.Request.Context(), ports.ListAllStudentsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Assigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		allEnrolledStudents, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), sid)
		available := filterAvailableUsers(allStudents, allEnrolledStudents)

		ctx.Header("Content-Type", "text/html")
		msg := "1 student assigned successfully"
		if len(studentIDs) != 1 {
			msg = fmt.Sprintf("%d students assigned successfully", len(studentIDs))
		}
		helpers.Toast(ctx, "Success", msg, toast.VariantSuccess)
		props := subjectStudentsTableProps(sid, students, available, totalCount, int32(limit), int32(offset), "")
		render.Render(ctx, components.UserTableRoot("students-user-table-root", props))
	}
}

func (h *UserHandler) SearchSubjectStudents() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
			return
		}

		search := strings.TrimSpace(ctx.PostForm("search"))
		if search == "" {
			search = strings.TrimSpace(ctx.Query("search"))
		}

		limitStr := ctx.DefaultQuery("limit", "12")
		offsetStr := ctx.DefaultQuery("offset", "0")
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		var (
			students   []domain.User
			totalCount int64
		)

		if search == "" {
			students, err = h.App.ListStudentsBySubjectIDPaginated(ctx.Request.Context(), subjectID, int32(limit), int32(offset))
			if err != nil {
				students = []domain.User{}
			}
			totalCount, err = h.App.CountStudentsBySubjectID(ctx.Request.Context(), subjectID)
			if err != nil {
				totalCount = 0
			}
		} else {
			students, err = h.App.SearchStudentsBySubjectID(ctx.Request.Context(), subjectID, search, int32(limit), int32(offset))
			if err != nil {
				students = []domain.User{}
			}
			totalCount, err = h.App.CountSearchStudentsBySubjectID(ctx.Request.Context(), subjectID, search)
			if err != nil {
				totalCount = 0
			}
		}

		allStudents, err := h.App.ListAllStudents(ctx.Request.Context(), ports.ListAllStudentsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			allStudents = []domain.User{}
		}
		allEnrolledStudents, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), subjectID)
		available := filterAvailableUsers(allStudents, allEnrolledStudents)

		ctx.Header("Content-Type", "text/html")
		props := subjectStudentsTableProps(subjectID, students, available, totalCount, int32(limit), int32(offset), search)
		render.Render(ctx, components.UserTableContainer(props))
	}
}

func subjectStudentsTableProps(subjectID uuid.UUID, students, allStudents []domain.User, totalCount int64, limit int32, offset int32, search string) components.UserTableProps {
	sid := subjectID.String()
	base := "/admin/dashboard/subject/" + sid
	baseURL := base + "/students/search"
	if search != "" {
		baseURL += "?search=" + url.QueryEscape(search)
	}
	return components.UserTableProps{
		Users:                   students,
		Title:                   "Students",
		BaseURL:                 baseURL,
		ID:                      "students-table",
		Search:                  true,
		SearchAPI:               base + "/students/search",
		AddUser:                 false,
		AddStudent:              true,
		Students:                allStudents,
		AddStudentAPI:           base + "/students/assign",
		StudentAssignSwapTarget: "#students-user-table-root",
		AddStudentCSVAPI:        base + "/students/assign-csv",
		ExportURL:               base + "/students/export",
		UnassignAPI:             base + "/students/unassign/:user_id",
		TotalCount:              totalCount,
		Limit:                   limit,
		Offset:                  offset,
	}
}

func (h *UserHandler) AssignStudentToSubjectCSV() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Student Failed", "Invalid subject ID", toast.VariantError)
			return
		}

		file, err := helpers.ParseCSVFile(ctx, "student_csv")
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Student Failed", "Invalid file: "+err.Error(), toast.VariantError)
			return
		}

		usernames, err := helpers.MapStudentCSVToStruct(file)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Student Failed", "Failed to parse CSV file", toast.VariantError)
			return
		}

		var studentIDs []uuid.UUID
		for _, username := range usernames {
			user, err := h.App.GetUserByUsername(ctx.Request.Context(), username)
			if err != nil {
				continue
			}
			if user.Role == domain.RoleStudent {
				studentIDs = append(studentIDs, user.ID)
			}
		}

		if len(studentIDs) == 0 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Student Failed", "No valid students found in CSV", toast.VariantError)
			return
		}

		for _, studentID := range studentIDs {
			err := h.App.AssignStudentsToSubject(ctx, subjectID, []uuid.UUID{studentID})
			if err != nil && !strings.Contains(err.Error(), "duplicate key value") && !strings.Contains(err.Error(), "unique constraint") {
				ctx.Header("HX-Reswap", "none")
				helpers.Toast(ctx, "Assign Students Failed", "Failed to assign student: "+err.Error(), toast.VariantError)
				return
			}
		}

		limitStr := ctx.DefaultQuery("limit", "12")
		offsetStr := ctx.DefaultQuery("offset", "0")
		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		students, err := h.App.ListStudentsBySubjectIDPaginated(ctx.Request.Context(), subjectID, int32(limit), int32(offset))
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Assigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		totalCount, _ := h.App.CountStudentsBySubjectID(ctx.Request.Context(), subjectID)
		allStudents, err := h.App.ListAllStudents(ctx.Request.Context(), ports.ListAllStudentsParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Assigned but failed to reload list: "+err.Error(), toast.VariantError)
			return
		}
		allEnrolledStudents, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), subjectID)
		available := filterAvailableUsers(allStudents, allEnrolledStudents)

		ctx.Header("Content-Type", "text/html")
		msg := "1 student assigned successfully from CSV"
		if len(studentIDs) != 1 {
			msg = fmt.Sprintf("%d students assigned successfully from CSV", len(studentIDs))
		}
		helpers.Toast(ctx, "Success", msg, toast.VariantSuccess)
		props := subjectStudentsTableProps(subjectID, students, available, totalCount, int32(limit), int32(offset), "")
		render.Render(ctx, components.UserTableRoot("students-user-table-root", props))
	}
}

func (h *UserHandler) EditUserPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/auth/login")
			return
		}

		adminUserClaim := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		targetUserIDStr := ctx.Param("id")
		targetUserID, err := uuid.Parse(targetUserIDStr)
		if err != nil {
			ctx.Redirect(302, "/admin/dashboard/users")
			return
		}

		targetUser, err := h.App.GetUserByID(ctx.Request.Context(), targetUserID)
		if err != nil {
			ctx.Redirect(302, "/admin/dashboard/users")
			return
		}

		var exams []domain.Exam
		var attempts []domain.ExamAttempt

		switch targetUser.Role {
		case domain.RoleStudent:
			exams, _ = h.App.ListExamsForStudent(ctx.Request.Context(), targetUserID)
			attempts, _ = h.App.ListAttemptsByStudent(ctx.Request.Context(), targetUserID)
		case domain.RoleInstructor, domain.RoleAdmin:
			exams, _ = h.App.ListExamsCreatedBy(ctx.Request.Context(), targetUserID)
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Edit User", page.EditUserPage(page.EditUserPageParam{
			AdminUser:  adminUserClaim,
			TargetUser: targetUser,
			Exams:      exams,
			Attempts:   attempts,
		})))
	}
}

func (h *UserHandler) EditUserInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		targetUserIDStr := ctx.Param("id")
		targetUserID, err := uuid.Parse(targetUserIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Failed", "Invalid user ID", toast.VariantError)
			return
		}

		username := ctx.PostForm("username")
		fullName := ctx.PostForm("full_name")
		role := ctx.PostForm("role")
		isActive := ctx.PostForm("is_active") == "true"

		if username == "" || fullName == "" || role == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Update Failed", "Please fill in all required fields", toast.VariantError)
			return
		}

		_, err = h.App.UpdateUserInfo(ctx.Request.Context(), ports.UpdateUserInfoParams{
			ID:       targetUserID,
			Username: username,
			FullName: fullName,
			Role:     domain.UserRole(role),
			IsActive: isActive,
		})

		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Update Failed", "Failed to update user: "+err.Error(), toast.VariantError)
			return
		}

		ctx.Header("HX-Reswap", "none")
		helpers.Toast(ctx, "Success", "User updated successfully", toast.VariantSuccess)
	}
}

// AdminSubjectTrackerWS handles real-time timer sync for the admin dashboard.
func (h *UserHandler) AdminSubjectTrackerWS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Status(400)
			return
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Re-calculate end time on every tick to pick up manual extensions/ending
			subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
			if err != nil || subject.Status != domain.SubjectStatusPublished {
				break
			}

			exams, err := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
			if err != nil {
				break
			}

			var examEndTime *time.Time
			for _, exam := range exams {
				if exam.Status == domain.ExamStatusPublished {
					t := exam.UpdatedAt.Add(time.Duration(subject.DurationMinutes) * time.Minute)
					if examEndTime == nil || t.After(*examEndTime) {
						examEndTime = &t
					}
				}
			}

			if examEndTime == nil {
				break
			}

			now := time.Now()
			remaining := int(examEndTime.Sub(now).Seconds())
			if remaining <= 0 {
				// Auto-submit all in-progress attempts for every published exam
				// and create+submit zero-score attempts for students who never started
				students, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), subjectID)

				for _, exam := range exams {
					if exam.Status != domain.ExamStatusPublished {
						continue
					}
					// 1) Submit in-progress attempts
					attempts, _ := h.App.ListInProgressAttemptsByExam(ctx.Request.Context(), exam.ID)
					questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), exam.ID)

					questionMarks := make(map[uuid.UUID]float64)
					for _, q := range questions {
						questionMarks[q.ID] = q.Marks
					}

					// Track which students already have an attempt
					studentsWithAttempt := make(map[uuid.UUID]bool)
					for _, attempt := range attempts {
						studentsWithAttempt[attempt.StudentID] = true
						answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
						var totalScore float64
						for _, answer := range answers {
							if answer.IsCorrect != nil && *answer.IsCorrect {
								if marks, ok := questionMarks[answer.QuestionID]; ok {
									totalScore += marks
								}
							}
						}
						_ = h.App.SubmitExamAttempt(ctx.Request.Context(), attempt.ID, totalScore)
					}

					// 2) Create + submit zero-score attempts for students who never started
					for _, student := range students {
						if studentsWithAttempt[student.ID] {
							continue
						}
						// Check if student already has a submitted attempt
						existingAttempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), exam.ID, student.ID)
						if err == nil && existingAttempt.Status != domain.AttemptStatusInProgress {
							continue // already submitted/graded
						}
						newAttempt, err := h.App.StartExamAttempt(ctx.Request.Context(), exam.ID, student.ID)
						if err != nil {
							continue
						}
						_ = h.App.SubmitExamAttempt(ctx.Request.Context(), newAttempt.ID, 0.0)
					}
				}

				_ = conn.WriteJSON(map[string]interface{}{
					"remaining_seconds": 0,
					"auto_submit":       true,
				})
				break
			} else {
				err := conn.WriteJSON(map[string]interface{}{
					"remaining_seconds": remaining,
					"auto_submit":       false,
				})
				if err != nil {
					break
				}
			}
		}
	}
}

func (h *UserHandler) UnassignStudentFromSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		studentID := ctx.Param("user_id")

		if subjectID == "" || studentID == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Student Failed", "Subject ID and Student ID are required", toast.VariantError)
			return
		}

		sid, err := uuid.Parse(subjectID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Student Failed", "Invalid subject ID", toast.VariantError)
			return
		}

		stid, err := uuid.Parse(studentID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Student Failed", "Invalid student ID", toast.VariantError)
			return
		}

		if err := h.App.UnassignStudentFromSubject(ctx, sid, stid); err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Unassign Student Failed", err.Error(), toast.VariantError)
			return
		}

		helpers.Toast(ctx, "Student Unassigned", "Successfully unassigned student from the subject", toast.VariantSuccess)

		// Re-render the table
		limit := int32(12)
		offset := int32(0)
		if l, err := strconv.Atoi(ctx.Query("limit")); err == nil {
			limit = int32(l)
		}
		if o, err := strconv.Atoi(ctx.Query("offset")); err == nil {
			offset = int32(o)
		}

		students, _ := h.App.ListStudentsBySubjectIDPaginated(ctx, sid, limit, offset)
		totalCount, _ := h.App.CountStudentsBySubjectID(ctx, sid)
		allStudents, _ := h.App.SearchStudents(ctx, ports.SearchStudentsParams{Search: "", Limit: 100, Offset: 0})

		ctx.Header("Content-Type", "text/html")
		components.UserTableContainer(components.UserTableProps{
			Users:                   students,
			Title:                   "Students",
			BaseURL:                 "/admin/dashboard/subject/" + subjectID + "/students/search",
			ID:                      "students-table",
			Search:                  true,
			SearchAPI:               "/admin/dashboard/subject/" + subjectID + "/students/search",
			TotalCount:              totalCount,
			Limit:                   limit,
			Offset:                  offset,
			AddStudent:              true,
			Students:                allStudents,
			AddStudentAPI:           "/admin/dashboard/subject/" + subjectID + "/students/assign",
			AddStudentCSVAPI:        "/admin/dashboard/subject/" + subjectID + "/students/assign-csv",
			StudentAssignSwapTarget: "#students-user-table-root",
			ExportURL:               "/admin/dashboard/subject/" + subjectID + "/students/export",
			UnassignAPI:             "/admin/dashboard/subject/" + subjectID + "/students/unassign/:user_id",
		}).Render(ctx, ctx.Writer)
	}
}

func (h *UserHandler) ExamAttemptsPDF() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exam ID"})
			return
		}

		exam, err := h.App.GetExamByID(ctx.Request.Context(), examID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), exam.SubjectID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
			return
		}

		attempts, err := h.App.ListAttemptsByExam(ctx.Request.Context(), examID)
		if err != nil {
			attempts = []domain.ExamAttempt{}
		}

		m := maroto.New(config.NewBuilder().Build())

		// Issue Date
		m.AddRows(
			row.New(5).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("Report Issue Date: %s", time.Now().Format("Jan 02, 2006")), props.Text{
						Size:  6,
						Align: align.Right,
						Color: &props.Color{Red: 148, Green: 163, Blue: 184},
					}),
				),
			),
		)

		// Header Bar
		m.AddRows(
			row.New(20).Add(
				col.New(4).Add(
					image.NewFromFile("./web/static/images/uotLogo.png", props.Rect{
						Center:  true,
						Percent: 90,
					}),
				),
				col.New(8).Add(
					text.New("EXAMINATION MANAGEMENT SYSTEM", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
						Align: align.Right,
						Top:   5,
						Color: &props.Color{Red: 30, Green: 58, Blue: 138}, // Dark Blue
					}),
				),
			),
			line.NewRow(1),
		)

		m.AddRows(row.New(10))

		// Title Section
		m.AddRows(
			row.New(15).Add(
				col.New(12).Add(
					text.New("STUDENT ATTEMPTS REPORT", props.Text{
						Size:  20,
						Style: fontstyle.Bold,
						Align: align.Center,
						Color: &props.Color{Red: 15, Green: 23, Blue: 42}, // Slate 900
					}),
				),
			),
			row.New(10).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("%s - %s", subject.Title, exam.Title), props.Text{
						Size:  12,
						Style: fontstyle.BoldItalic,
						Align: align.Center,
						Color: &props.Color{Red: 71, Green: 85, Blue: 105}, // Slate 600
					}),
				),
			),
		)

		m.AddRows(row.New(10))

		// Exam Info Box
		m.AddRows(
			row.New(8).Add(
				col.New(12).Add(
					text.New("EXAM INFORMATION", props.Text{Size: 7, Style: fontstyle.Bold, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Left: 3, Top: 2}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
			row.New(18).Add(
				col.New(3).Add(
					text.New("Total Marks", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Left: 3}),
					text.New(fmt.Sprintf("%.2f", exam.TotalMarks), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4, Left: 3}),
				),
				col.New(3).Add(
					text.New("Pass Score", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(fmt.Sprintf("%.2f", exam.PassScore), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
				col.New(3).Add(
					text.New("Total Students", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(fmt.Sprintf("%d", len(attempts)), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
				col.New(3).Add(
					text.New("Report Date", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(time.Now().Format("Jan 02, 2006"), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
		)

		m.AddRows(row.New(15))

		// Attempts Table Header
		m.AddRows(
			row.New(10).Add(
				col.New(4).Add(text.New("STUDENT NAME", props.Text{Size: 9, Style: fontstyle.Bold, Color: &props.Color{Red: 255, Green: 255, Blue: 255}, Left: 3})),
				col.New(2).Add(text.New("ID", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(2).Add(text.New("STATUS", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(2).Add(text.New("SCORE", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(2).Add(text.New("DATE", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}, Right: 3})),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 30, Green: 41, Blue: 59}}), // Slate 800
		)

		for i, attempt := range attempts {
			score := "-"
			if attempt.Score != nil {
				score = fmt.Sprintf("%.2f", *attempt.Score)
			}

			statusColor := &props.Color{Red: 71, Green: 85, Blue: 105} // Slate 600
			if attempt.Status == domain.AttemptStatusSubmitted || attempt.Status == domain.AttemptStatusGraded {
				if attempt.Score != nil && *attempt.Score >= exam.PassScore {
					statusColor = &props.Color{Red: 21, Green: 128, Blue: 61} // Green 700
				} else {
					statusColor = &props.Color{Red: 185, Green: 28, Blue: 28} // Red 700
				}
			}

			rowStyle := &props.Cell{BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255}}
			if i%2 == 1 {
				rowStyle.BackgroundColor = &props.Color{Red: 241, Green: 245, Blue: 249} // Slate 100
			}

			m.AddRows(
				row.New(10).Add(
					col.New(4).Add(text.New(attempt.StudentName, props.Text{Size: 9, Top: 2, Left: 3})),
					col.New(2).Add(text.New(attempt.StudentUsername, props.Text{Size: 8, Align: align.Center, Top: 2})),
					col.New(2).Add(text.New(string(attempt.Status), props.Text{Size: 8, Align: align.Center, Color: statusColor, Top: 2, Style: fontstyle.Bold})),
					col.New(2).Add(text.New(score, props.Text{Size: 9, Align: align.Center, Top: 2, Style: fontstyle.Bold})),
					col.New(2).Add(text.New(attempt.StartedAt.Format("Jan 02, 2006"), props.Text{Size: 8, Align: align.Center, Top: 2, Right: 3})),
				).WithStyle(rowStyle),
			)
		}

		m.AddRows(row.New(20))

		// Footer
		m.AddRows(
			row.New(10).Add(
				col.New(12).Add(
					text.New("This is a computer-generated report. All records are stored securely in the system.", props.Text{
						Size:  8,
						Align: align.Center,
						Style: fontstyle.Italic,
						Color: &props.Color{Red: 100, Green: 116, Blue: 139},
					}),
				),
			),
		)

		document, _ := m.Generate()

		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=exam_attempts_%s.pdf", url.PathEscape(exam.Title)))
		ctx.Header("Content-Type", "application/pdf")
		ctx.Data(200, "application/pdf", document.GetBytes())
	}
}
func (h *UserHandler) AdminSettingsPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)

		// Fetch settings from DB
		var systemName, adminEmail string
		_ = h.App.DB().Pool.QueryRow(ctx.Request.Context(), "SELECT value FROM system_settings WHERE key = 'system_name'").Scan(&systemName)
		_ = h.App.DB().Pool.QueryRow(ctx.Request.Context(), "SELECT value FROM system_settings WHERE key = 'admin_email'").Scan(&adminEmail)

		if systemName == "" {
			systemName = "UoT Examination System"
		}
		if adminEmail == "" {
			adminEmail = "admin@uot.edu.iq"
		}

		params := page.SettingsPageParam{
			Username:   username,
			FullName:   fullname,
			UserID:     userID.String(),
			SystemName: systemName,
			AdminEmail: adminEmail,
			AppVersion: h.viperConfig.AppVersion,
			GoVersion:  runtime.Version(),
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			CPUs:       runtime.NumCPU(),
			Uptime:     time.Since(startTime).Truncate(time.Second).String(),
		}
		render.Render(ctx, pages.BasePage("System Settings", page.SettingsPage(params)))
	}
}

func (h *UserHandler) UpdateSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		systemName := ctx.PostForm("system_name")
		adminEmail := ctx.PostForm("admin_email")

		if systemName != "" {
			_, _ = h.App.DB().Pool.Exec(ctx.Request.Context(),
				"INSERT INTO system_settings (key, value) VALUES ('system_name', $1) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value",
				systemName)
		}
		if adminEmail != "" {
			_, _ = h.App.DB().Pool.Exec(ctx.Request.Context(),
				"INSERT INTO system_settings (key, value) VALUES ('admin_email', $1) ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value",
				adminEmail)
		}

		helpers.Toast(ctx, "Settings Saved", "System configurations have been updated successfully.", toast.VariantSuccess)
	}
}

func (h *UserHandler) ExportBackup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCfg := h.viperConfig.Database

		// Set PGPASSWORD environment variable for pg_dump
		cmd := exec.Command("pg_dump",
			"-h", dbCfg.Host,
			"-p", strconv.Itoa(dbCfg.Port),
			"-U", dbCfg.Username,
			"-d", dbCfg.DatabaseName,
			"--no-owner", "--no-privileges",
		)
		cmd.Env = append(cmd.Env, "PGPASSWORD="+dbCfg.Password)

		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=backup_%s.sql", time.Now().Format("2006-01-02_150405")))
		ctx.Header("Content-Type", "application/sql")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if err := cmd.Start(); err != nil {
			h.logger.LogErrorWithLevel("error", "BACKUP_FAILED", "BACKUP_FAILED", "Failed to start pg_dump", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// Stream the output directly to the response writer
		_, _ = io.Copy(ctx.Writer, stdout)

		if err := cmd.Wait(); err != nil {
			h.logger.LogErrorWithLevel("error", "BACKUP_FAILED", "BACKUP_FAILED", "pg_dump failed during execution", err)
			// Headers already sent, so we can't change the status code here
			return
		}
	}
}

func (h *UserHandler) SubjectOverallAttemptsPDF() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subject ID"})
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
			return
		}

		attempts, err := h.App.ListOverallAttemptsBySubject(ctx.Request.Context(), subjectID)
		if err != nil {
			attempts = []domain.SubjectAttempt{}
		}

		m := maroto.New(config.NewBuilder().Build())

		// Issue Date
		m.AddRows(
			row.New(5).Add(
				col.New(12).Add(
					text.New(fmt.Sprintf("Report Issue Date: %s", time.Now().Format("Jan 02, 2006")), props.Text{
						Size:  6,
						Align: align.Right,
						Color: &props.Color{Red: 148, Green: 163, Blue: 184},
					}),
				),
			),
		)

		// Header Bar
		m.AddRows(
			row.New(20).Add(
				col.New(4).Add(
					image.NewFromFile("./web/static/images/uotLogo.png", props.Rect{
						Center:  true,
						Percent: 90,
					}),
				),
				col.New(8).Add(
					text.New("EXAMINATION MANAGEMENT SYSTEM", props.Text{
						Size:  14,
						Style: fontstyle.Bold,
						Align: align.Right,
						Top:   5,
						Color: &props.Color{Red: 30, Green: 58, Blue: 138}, // Dark Blue
					}),
				),
			),
			line.NewRow(1),
		)

		m.AddRows(row.New(10))

		// Title Section
		m.AddRows(
			row.New(15).Add(
				col.New(12).Add(
					text.New("OVERALL PERFORMANCE REPORT", props.Text{
						Size:  20,
						Style: fontstyle.Bold,
						Align: align.Center,
						Color: &props.Color{Red: 15, Green: 23, Blue: 42}, // Slate 900
					}),
				),
			),
			row.New(10).Add(
				col.New(12).Add(
					text.New(subject.Title, props.Text{
						Size:  12,
						Style: fontstyle.BoldItalic,
						Align: align.Center,
						Color: &props.Color{Red: 71, Green: 85, Blue: 105}, // Slate 600
					}),
				),
			),
		)

		m.AddRows(row.New(10))

		// Subject Info Box
		m.AddRows(
			row.New(8).Add(
				col.New(12).Add(
					text.New("SUBJECT INFORMATION", props.Text{Size: 7, Style: fontstyle.Bold, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Left: 3, Top: 2}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
			row.New(18).Add(
				col.New(3).Add(
					text.New("Total Marks", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Left: 3}),
					text.New(fmt.Sprintf("%.2f", subject.TotalMarks), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4, Left: 3}),
				),
				col.New(3).Add(
					text.New("Pass Score", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(fmt.Sprintf("%.2f", subject.PassScore), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
				col.New(3).Add(
					text.New("Total Students", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(fmt.Sprintf("%d", len(attempts)), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
				col.New(3).Add(
					text.New("Report Date", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(time.Now().Format("Jan 02, 2006"), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
		)

		m.AddRows(row.New(15))

		// Performance Table Header
		m.AddRows(
			row.New(10).Add(
				col.New(3).Add(text.New("STUDENT NAME", props.Text{Size: 9, Style: fontstyle.Bold, Color: &props.Color{Red: 255, Green: 255, Blue: 255}, Left: 3})),
				col.New(1).Add(text.New("ID", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(2).Add(text.New("DURATION", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(2).Add(text.New("SCORE", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(1).Add(text.New("RESULT", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(3).Add(text.New("SUBMITTED AT", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}, Right: 3})),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 30, Green: 41, Blue: 59}}), // Slate 800
		)

		for i, attempt := range attempts {
			result := "FAIL"
			resultColor := &props.Color{Red: 185, Green: 28, Blue: 28} // Red 700
			if attempt.TotalScore >= subject.PassScore {
				result = "PASS"
				resultColor = &props.Color{Red: 21, Green: 128, Blue: 61} // Green 700
			}

			submittedAt := "-"
			if attempt.LastSubmittedAt != nil {
				submittedAt = attempt.LastSubmittedAt.Format("2006-01-02 15:04")
			}

			rowStyle := &props.Cell{BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255}}
			if i%2 == 1 {
				rowStyle.BackgroundColor = &props.Color{Red: 241, Green: 245, Blue: 249} // Slate 100
			}

			m.AddRows(
				row.New(10).Add(
					col.New(3).Add(text.New(attempt.StudentName, props.Text{Size: 9, Top: 2, Left: 3})),
					col.New(1).Add(text.New(attempt.StudentUsername, props.Text{Size: 8, Align: align.Center, Top: 2})),
					col.New(2).Add(text.New(formatDuration(attempt.TotalDuration), props.Text{Size: 8, Align: align.Center, Top: 2})),
					col.New(2).Add(text.New(fmt.Sprintf("%.2f", attempt.TotalScore), props.Text{Size: 9, Align: align.Center, Top: 2, Style: fontstyle.Bold})),
					col.New(1).Add(text.New(result, props.Text{Size: 8, Align: align.Center, Color: resultColor, Top: 2, Style: fontstyle.Bold})),
					col.New(3).Add(text.New(submittedAt, props.Text{Size: 8, Align: align.Center, Top: 2, Right: 3})),
				).WithStyle(rowStyle),
			)
		}

		m.AddRows(row.New(40))

		// Signature Section
		m.AddRows(
			row.New(20).Add(
				col.New(7).Add(text.New("")),
				col.New(5).Add(
					line.New(props.Line{Thickness: 0.5, Color: &props.Color{Red: 148, Green: 163, Blue: 184}}),
					text.New("Examiner Official Signature", props.Text{Size: 8, Align: align.Center, Top: 4, Color: &props.Color{Red: 71, Green: 85, Blue: 105}}),
				),
			),
		)
		// Footer
		m.AddRows(
			row.New(10).Add(
				col.New(12).Add(
					text.New("This is a computer-generated report. All records are stored securely in the system.", props.Text{
						Size:  8,
						Align: align.Center,
						Style: fontstyle.Italic,
						Color: &props.Color{Red: 100, Green: 116, Blue: 139},
					}),
				),
			),
		)

		document, _ := m.Generate()

		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=overall_performance_%s.pdf", url.PathEscape(subject.Title)))
		ctx.Header("Content-Type", "application/pdf")
		ctx.Data(200, "application/pdf", document.GetBytes())
	}
}

func formatDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func (h *UserHandler) AdminSubjectEndTimer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid subject ID", toast.VariantError)
			return
		}

		// Check if timer is already expired
		subject, _ := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		exams, _ := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)

		var examEndTime *time.Time
		for _, exam := range exams {
			if exam.Status == domain.ExamStatusPublished {
				t := exam.UpdatedAt.Add(time.Duration(subject.DurationMinutes) * time.Minute)
				if examEndTime == nil || t.After(*examEndTime) {
					examEndTime = &t
				}
			}
		}

		if examEndTime == nil || time.Now().After(*examEndTime) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Notice", "The examination session has already concluded.", toast.VariantInfo)
			return
		}

		err = h.App.EndExamsTimer(ctx.Request.Context(), subjectID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Failed to end timer: "+err.Error(), toast.VariantError)
			return
		}

		ctx.Header("HX-Reswap", "none")
		helpers.Toast(ctx, "Success", "Exam timer ended. Students will be auto-submitted shortly.", toast.VariantSuccess)
	}
}

func (h *UserHandler) AdminSubjectExtendTimer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid subject ID", toast.VariantError)
			return
		}

		minutesStr := ctx.PostForm("minutes")
		var minutes int32
		fmt.Sscanf(minutesStr, "%d", &minutes)
		if minutes <= 0 {
			minutes = 10 // Default extension
		}

		// Check if timer is already expired
		subject, _ := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		exams, _ := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)

		var examEndTime *time.Time
		for _, exam := range exams {
			if exam.Status == domain.ExamStatusPublished {
				t := exam.UpdatedAt.Add(time.Duration(subject.DurationMinutes) * time.Minute)
				if examEndTime == nil || t.After(*examEndTime) {
					examEndTime = &t
				}
			}
		}

		if examEndTime == nil || time.Now().After(*examEndTime) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Notice", "The examination session has already concluded.", toast.VariantInfo)
			return
		}

		err = h.App.ShiftExamsTimer(ctx.Request.Context(), subjectID, minutes)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Failed to extend timer: "+err.Error(), toast.VariantError)
			return
		}

		ctx.Header("HX-Reswap", "none")
		helpers.Toast(ctx, "Success", fmt.Sprintf("Timer extended by %d minutes", minutes), toast.VariantSuccess)
	}
}
