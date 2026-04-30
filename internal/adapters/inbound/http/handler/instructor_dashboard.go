package handler

import (
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/pages"
	instructorPages "uot-exam/web/templates/pages/instructor_dashboard"
	"uot-exam/web/templates/render"
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/web/templates/components/toast"
	adminComponents "uot-exam/web/templates/pages/admin_dashboard/components"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

// InstructorDashboardRender renders the instructor dashboard showing assigned subjects.
func (h *UserHandler) InstructorDashboardRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		instructorUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		subjects, err := h.App.ListSubjectsForInstructor(ctx.Request.Context(), userID)
		if err != nil {
			subjects = []domain.Subject{}
		}

		var subjectInfos []instructorPages.SubjectInfo
		for _, subject := range subjects {
			totalExams, _ := h.App.CountExamsBySubject(ctx.Request.Context(), subject.ID)
			subjectInfos = append(subjectInfos, instructorPages.SubjectInfo{
				Subject:    subject,
				TotalExams: totalExams,
			})
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Instructor Dashboard", instructorPages.DashboardPage(instructorPages.DashboardPageParam{
			Instructor: instructorUser,
			Subjects:   subjectInfos,
		})))
	}
}

// InstructorSubjectView renders the subject detail page for instructors.
func (h *UserHandler) InstructorSubjectView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		instructorUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		// Verify that the instructor is assigned to this subject
		assignedInstructors, err := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), subjectID)
		isAssigned := false
		for _, inst := range assignedInstructors {
			if inst.ID == userID {
				isAssigned = true
				break
			}
		}

		if !isAssigned {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		if err != nil {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		exams, err := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
		if err != nil {
			exams = []domain.Exam{}
		}

		students, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), subjectID)
		allStudents, _ := h.App.SearchStudents(ctx.Request.Context(), ports.SearchStudentsParams{Search: "", Limit: 100, Offset: 0})

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Subject - "+subject.Title, instructorPages.SubjectViewPage(instructorPages.SubjectViewParam{
			Instructor: instructorUser,
			Subject:    subject,
			Exams:       exams,
			Students:    students,
			AllStudents: allStudents,
		})))
	}
}

// InstructorEditExamPageRender renders the exam editing page for instructors.
func (h *UserHandler) InstructorEditExamPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		instructorUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		exam, err := h.App.GetExamByID(ctx.Request.Context(), examID)
		if err != nil {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), exam.SubjectID)
		if err != nil {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		// Verify that the instructor is assigned to this subject
		assignedInstructors, _ := h.App.ListInstructorsBySubjectID(ctx.Request.Context(), subject.ID)
		isAssigned := false
		for _, inst := range assignedInstructors {
			if inst.ID == userID {
				isAssigned = true
				break
			}
		}

		if !isAssigned {
			ctx.Redirect(302, "/instructor/dashboard")
			return
		}

		questions, err := h.App.ListQuestionsByExam(ctx.Request.Context(), examID)
		if err != nil {
			questions = []domain.Question{}
		}

		for i := range questions {
			choices, _ := h.App.ListChoicesByQuestion(ctx.Request.Context(), questions[i].ID)
			questions[i].Choices = choices
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Edit Exam - "+exam.Title, instructorPages.EditExamPage(instructorPages.EditExamPageParam{
			Instructor: instructorUser,
			Subject:    subject,
			Exam:       exam,
			Questions:  questions,
		})))
	}
}

// SearchInstructorSubjects handles searching subjects for an instructor.
func (h *UserHandler) SearchInstructorSubjects() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)
		if userID == uuid.Nil {
			ctx.Status(401)
			return
		}

		search := ctx.PostForm("search")
		subjects, err := h.App.ListSubjectsForInstructor(ctx.Request.Context(), userID)
		if err != nil {
			subjects = []domain.Subject{}
		}

		var subjectInfos []instructorPages.SubjectInfo
		for _, subject := range subjects {
			if search != "" && !strings.Contains(strings.ToLower(subject.Title), strings.ToLower(search)) {
				continue
			}
			totalExams, _ := h.App.CountExamsBySubject(ctx.Request.Context(), subject.ID)
			subjectInfos = append(subjectInfos, instructorPages.SubjectInfo{
				Subject:    subject,
				TotalExams: totalExams,
			})
		}

		ctx.Header("Content-Type", "text/html")
		instructorPages.SubjectGrid(subjectInfos).Render(ctx.Request.Context(), ctx.Writer)
	}
}

// SearchSubjectStudentsInstructor handles searching students for a subject in the instructor dashboard.
func (h *UserHandler) SearchSubjectStudentsInstructor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Status(400)
			return
		}

		search := ctx.PostForm("search")
		students, err := h.App.SearchStudentsBySubjectID(ctx.Request.Context(), subjectID, search, 100, 0)
		if err != nil {
			students = []domain.User{}
		}

		allStudents, _ := h.App.SearchStudents(ctx.Request.Context(), ports.SearchStudentsParams{Search: "", Limit: 100, Offset: 0})

		ctx.Header("Content-Type", "text/html")
		adminComponents.UserTableContainer(adminComponents.UserTableProps{
			Users:     students,
			Title:     "Students",
			BaseURL:   "/instructor/subject/" + subjectIDStr + "/students/search",
			ID:        "students-table",
			Search:    true,
			SearchAPI: "/instructor/subject/" + subjectIDStr + "/students/search",
			AddStudent: true,
			Students: allStudents,
			AddStudentAPI: "/instructor/subject/" + subjectIDStr + "/students/assign",
			StudentAssignSwapTarget: "#students-user-table-root",
            HideActions: true,
		}).Render(ctx.Request.Context(), ctx.Writer)
	}
}

// AssignStudentToSubjectInstructor handles assigning students to a subject in the instructor dashboard.
func (h *UserHandler) AssignStudentToSubjectInstructor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectIDStr := ctx.Param("id")
		rawIDs := ctx.PostForm("student_id")
		
		if subjectIDStr == "" || strings.TrimSpace(rawIDs) == "" {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", "Invalid input", toast.VariantError)
			return
		}

		subjectID, err := uuid.Parse(subjectIDStr)
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

		err = h.App.AssignStudentsToSubject(ctx.Request.Context(), subjectID, studentIDs)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Assign Students Failed", err.Error(), toast.VariantError)
			return
		}

		// Success toast
		helpers.Toast(ctx, "Students Assigned", "Successfully assigned students to the subject", toast.VariantSuccess)

		// Re-render the entire root to refresh the list
		students, _ := h.App.ListStudentsBySubjectID(ctx.Request.Context(), subjectID)
		allStudents, _ := h.App.SearchStudents(ctx.Request.Context(), ports.SearchStudentsParams{Search: "", Limit: 100, Offset: 0})

		ctx.Header("Content-Type", "text/html")
		adminComponents.UserTableRoot("students-user-table-root", adminComponents.UserTableProps{
			Users:     students,
			Title:     "Students",
			BaseURL:   "/instructor/subject/" + subjectIDStr + "/students/search",
			ID:        "students-table",
			Search:    true,
			SearchAPI: "/instructor/subject/" + subjectIDStr + "/students/search",
			AddStudent: true,
			Students: allStudents,
			AddStudentAPI: "/instructor/subject/" + subjectIDStr + "/students/assign",
			StudentAssignSwapTarget: "#students-user-table-root",
            HideActions: true,
		}).Render(ctx.Request.Context(), ctx.Writer)
	}
}
