package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/components/toast"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/components"
	"uot-exam/web/templates/pages/admin_dashboard/page"
	"uot-exam/web/templates/render"

	"uot-exam/internal/adapters/inbound/http/helpers"

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

		username, fullname, _ := parseUsername(ctx)
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

		username, fullname, _ := parseUsername(ctx)

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
		limitStr := ctx.DefaultQuery("limit", "12")
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

		username, fullname, _ := parseUsername(ctx)

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

func parseUsername(ctx *gin.Context) (string, string, uuid.UUID) {
	usernameVal, exists := ctx.Get("username")
	if !exists {
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
		return "", "", uuid.Nil
	}
	username := usernameVal.(string)
	fullnameVal, exists := ctx.Get("fullname")
	if !exists {
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
		fmt.Println("fullname not found")
		return "", "", uuid.Nil
	}
	fullname := fullnameVal.(string)
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
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

		username, fullname, _ := parseUsername(ctx)

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

		if ctx.GetHeader("HX-Target") == "toast-container" {
			ctx.Header("HX-Redirect", "/admin/dashboard/exams")
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

		username, fullname, _ := parseUsername(ctx)

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
		ctx.Redirect(http.StatusSeeOther, "/admin/login")
	}
}

func (h *UserHandler) EditSubjectPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID := ctx.Param("id")
		username, fullname, userID := parseUsername(ctx)
		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), uuid.MustParse(subjectID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		instructors, err := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), uuid.MustParse(subjectID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		exams, err := h.App.ListExamsBySubject(ctx.Request.Context(), uuid.MustParse(subjectID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		render.Render(ctx, pages.BasePage("Edit Subject", page.EditSubjectPage(page.EditSubjectPageParam{
			Subject:     subject,
			Username:    username,
			FullName:    fullname,
			Instructors: instructors,
			Exams:       exams,
			UserID:      userID,
		})))
	}
}

func (h *UserHandler) CreateExam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		title := ctx.PostForm("title")
		subjectIDStr := ctx.PostForm("subject_id")
		createdByStr := ctx.PostForm("created_by")
		description := ctx.PostForm("description")
		durationStr := ctx.PostForm("duration")
		totalMarksStr := ctx.PostForm("total_marks")
		passScoreStr := ctx.PostForm("pass_marks")

		if helpers.IsTrimmedEmpty(title) || helpers.IsTrimmedEmpty(subjectIDStr) || helpers.IsTrimmedEmpty(createdByStr) {
			helpers.Toast(ctx, "Create Exam Failed", "Missing required fields", toast.VariantError)
			return
		}

		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			helpers.Toast(ctx, "Create Exam Failed", "Invalid subject ID", toast.VariantError)
			return
		}

		createdBy, err := uuid.Parse(createdByStr)
		if err != nil {
			helpers.Toast(ctx, "Create Exam Failed", "Invalid creator ID", toast.VariantError)
			return
		}

		// Verify if the user exists (prevents foreign key violation due to stale sessions)
		if _, err := h.App.GetUserByID(ctx, createdBy); err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Session expired or user not found. Please log out and log in again.", toast.VariantError)
			return
		}

		durationParts := strings.Split(durationStr, ":")
		if len(durationParts) != 2 {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid duration format (expected HH:MM)", toast.VariantError)
			return
		}
		hrs, err1 := strconv.Atoi(durationParts[0])
		m, err2 := strconv.Atoi(durationParts[1])
		if err1 != nil || err2 != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid duration values", toast.VariantError)
			return
		}
		duration := hrs*60 + m
		totalMarks, err := strconv.Atoi(totalMarksStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid total marks", toast.VariantError)
			return
		}
		passScore, err := strconv.Atoi(passScoreStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Exam Failed", "Invalid pass score", toast.VariantError)
			return
		}

		_, err = h.App.CreateExam(ctx, ports.CreateExamParams{
			Title:           title,
			SubjectID:       subjectID,
			CreatedBy:       createdBy,
			Description:     &description,
			DurationMinutes: int32(duration),
			TotalMarks:      int32(totalMarks),
			PassScore:       int32(passScore),
			Status:          domain.ExamStatusDraft,
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

		_, err := h.App.UpdateSubject(ctx, domain.Subject{
			ID:          uuid.MustParse(id),
			Title:       name,
			Description: &description,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Subject Failed", "Failed to update subject", toast.VariantError)
			return
		}

		helpers.Toast(ctx, "Edit Subject Success", "Subject updated successfully", toast.VariantSuccess)
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

		_, err := h.App.CreateSubject(ctx, domain.Subject{
			Title:       name,
			Description: &description,
		})
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Subject Failed", "Failed to create subject", toast.VariantError)
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

		username, fullname, _ := parseUsername(ctx)
		render.Render(ctx, pages.BasePage("Edit Exam", page.EditExamPage(page.EditExamPageParam{
			Exam:      exam,
			Subject:   subject,
			Username:  username,
			FullName:  fullname,
			Questions: []domain.Question{},
			Choices:   []domain.Choice{},
		})))
	}
}

func (h *UserHandler) EditExamInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		name := ctx.PostForm("exam_name")
		description := ctx.PostForm("exam_description")
		durationStr := ctx.PostForm("exam_duration")
		passScoreStr := ctx.PostForm("exam_pass_score")
		totalMarksStr := ctx.PostForm("exam_total_marks")
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

		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			duration = 0
		}
		passScore, err := strconv.Atoi(passScoreStr)
		if err != nil {
			passScore = 0
		}
		totalMarks, err := strconv.Atoi(totalMarksStr)
		if err != nil {
			totalMarks = 0
		}

		_, err = h.App.UpdateExam(ctx, ports.UpdateExamParams{
			ID:              examID,
			Title:           name,
			Description:     &description,
			DurationMinutes: int32(duration),
			PassScore:       int32(passScore),
			TotalMarks:      int32(totalMarks),
			Status:          domain.ExamStatus(statusStr),
		})

		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Edit Exam Failed", "Failed to update exam: "+err.Error(), toast.VariantError)
			return
		}

		helpers.Toast(ctx, "Edit Exam Success", "Exam updated successfully", toast.VariantSuccess)
	}
}
