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
		durationStr := ctx.PostForm("duration")
		totalMarksStr := ctx.PostForm("total_marks")
		passScoreStr := ctx.PostForm("pass_marks")

		if helpers.IsTrimmedEmpty(title) || helpers.IsTrimmedEmpty(subjectIDStr) {
			helpers.Toast(ctx, "Create Exam Failed", "Missing required fields", toast.VariantError)
			return
		}

		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			helpers.Toast(ctx, "Create Exam Failed", "Invalid subject ID", toast.VariantError)
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

		questions, err := h.App.ListQuestionsByExam(ctx.Request.Context(), examIDUUID)
		if err != nil {
			questions = []domain.Question{}
		}

		var choices []domain.Choice

		for i := range questions {
			choices, err = h.App.ListChoicesByQuestion(ctx.Request.Context(), questions[i].ID)
			if err != nil {
				choices = []domain.Choice{}
			}
			questions[i].Choices = choices
		}

		username, fullname, _ := parseUsername(ctx)
		render.Render(ctx, pages.BasePage("Edit Exam", page.EditExamPage(page.EditExamPageParam{
			Exam:      exam,
			Subject:   subject,
			Username:  username,
			FullName:  fullname,
			Questions: questions,
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
		questionMarks := ctx.PostForm("question_marks")
		correctChoice := ctx.PostForm("correct_choice")
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

		if helpers.IsTrimmedEmpty(questionTitle) || helpers.IsTrimmedEmpty(questionText) ||
			helpers.IsTrimmedEmpty(questionType) || helpers.IsTrimmedEmpty(questionMarks) || len(choices) == 0 ||
			helpers.IsTrimmedEmpty(correctChoice) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "All fields are required", toast.VariantError)
			return
		}

		questionMarksInt, err := strconv.Atoi(questionMarks)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Create Question Failed", "Invalid question marks", toast.VariantError)
			return
		}

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
			Marks:         questionMarksInt,
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
			Marks:         questionMarksInt,
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

		file, examID, parsedUUID, err := helpers.ParseCSVFile(ctx)
		if err != nil {
			return
		}

		questions, err := helpers.MapCSVToStruct(file)
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
