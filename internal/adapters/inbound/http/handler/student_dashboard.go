package handler

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/url"
	"time"

	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/domain"
	"uot-exam/web/templates/components/toast"
	"uot-exam/web/templates/pages"
	studentPages "uot-exam/web/templates/pages/student_dashboard"
	"uot-exam/web/templates/render"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

// StudentDashboardRender renders the student dashboard showing enrolled subjects.
func (h *UserHandler) StudentDashboardRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		studentUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		subjects, err := h.App.ListSubjectsForStudent(ctx.Request.Context(), userID)
		if err != nil {
			subjects = []domain.Subject{}
		}

		// Build subject info with exam counts
		type SubjectInfo struct {
			Subject        domain.Subject
			TotalExams     int64
			SubmittedExams int64
		}

		var subjectInfos []studentPages.SubjectInfo
		for _, subject := range subjects {
			totalExams, _ := h.App.CountPublishedExamsBySubject(ctx.Request.Context(), subject.ID)
			submittedExams, _ := h.App.CountSubmittedAttemptsBySubjectForStudent(ctx.Request.Context(), subject.ID, userID)
			subjectInfos = append(subjectInfos, studentPages.SubjectInfo{
				Subject:        subject,
				TotalExams:     totalExams,
				SubmittedExams: submittedExams,
			})
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Student Dashboard", studentPages.DashboardPage(studentPages.DashboardPageParam{
			Student:  studentUser,
			Subjects: subjectInfos,
		})))
	}
}

// StudentSubjectView renders the subject detail page showing its exams and their statuses.
func (h *UserHandler) StudentSubjectView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		studentUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		exams, err := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
		if err != nil {
			exams = []domain.Exam{}
		}

		// Filter only published exams for students
		var publishedExams []domain.Exam
		var examEndTime *time.Time
		for _, exam := range exams {
			if exam.Status == domain.ExamStatusPublished {
				publishedExams = append(publishedExams, exam)
				t := exam.UpdatedAt.Add(time.Duration(subject.DurationMinutes) * time.Minute)
				if examEndTime == nil || t.After(*examEndTime) {
					examEndTime = &t
				}
			}
		}

		// Get attempt status for each exam
		var examInfos []studentPages.ExamInfo
		for _, exam := range publishedExams {
			attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), exam.ID, userID)
			info := studentPages.ExamInfo{
				Exam: exam,
			}
			if err == nil {
				info.Attempt = &attempt
			}
			examInfos = append(examInfos, info)
		}

		// Check if all exams submitted to reveal scores
		totalExams, _ := h.App.CountPublishedExamsBySubject(ctx.Request.Context(), subjectID)
		submittedExams, _ := h.App.CountSubmittedAttemptsBySubjectForStudent(ctx.Request.Context(), subjectID, userID)
		allSubmitted := totalExams > 0 && submittedExams >= totalExams

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Subject - "+subject.Title, studentPages.SubjectViewPage(studentPages.SubjectViewParam{
			Student:      studentUser,
			Subject:      subject,
			Exams:        examInfos,
			AllSubmitted: allSubmitted,
			ExamEndTime:  examEndTime,
		})))
	}
}

// StudentStartExam starts or resumes an exam attempt.
func (h *UserHandler) StudentStartExam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)

		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid exam ID", toast.VariantError)
			return
		}

		// Check for existing in-progress attempt
		attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), examID, userID)
		if err == nil && attempt.Status == domain.AttemptStatusInProgress {
			// Resume existing attempt
			ctx.Header("HX-Redirect", fmt.Sprintf("/student/exam/%s/take", examID.String()))
			return
		}

		if err == nil && (attempt.Status == domain.AttemptStatusSubmitted || attempt.Status == domain.AttemptStatusGraded) {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Already Submitted", "You have already submitted this exam", toast.VariantError)
			return
		}

		// Start new attempt
		_, err = h.App.StartExamAttempt(ctx.Request.Context(), examID, userID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Failed to start exam: "+err.Error(), toast.VariantError)
			return
		}

		ctx.Header("HX-Redirect", fmt.Sprintf("/student/exam/%s/take", examID.String()))
	}
}

// StudentExamView renders the exam-taking page.
func (h *UserHandler) StudentExamView() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		studentUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		exam, err := h.App.GetExamByID(ctx.Request.Context(), examID)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), exam.SubjectID)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), examID, userID)
		if err != nil {
			// No attempt found, redirect to subject page
			ctx.Redirect(302, fmt.Sprintf("/student/subject/%s", exam.SubjectID.String()))
			return
		}

		if attempt.Status != domain.AttemptStatusInProgress {
			// Attempt already submitted
			ctx.Redirect(302, fmt.Sprintf("/student/subject/%s", exam.SubjectID.String()))
			return
		}

		questions, err := h.App.ListQuestionsByExam(ctx.Request.Context(), examID)
		if err != nil {
			questions = []domain.Question{}
		}

		// Load choices for each question
		for i, question := range questions {
			choices, err := h.App.ListChoicesByQuestion(ctx.Request.Context(), question.ID)
			if err == nil {
				questions[i].Choices = choices
			}
		}

		shuffleQuestionsAndChoices(userID, examID, questions)

		// Get existing answers
		answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
		answeredMap := make(map[uuid.UUID]uuid.UUID)
		for _, answer := range answers {
			answeredMap[answer.QuestionID] = answer.SelectedChoiceID
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Exam - "+exam.Title, studentPages.ExamTakePage(studentPages.ExamTakeParam{
			Student:         studentUser,
			Exam:            exam,
			Subject:         subject,
			Attempt:         attempt,
			Questions:       questions,
			AnsweredChoices: answeredMap,
		})))
	}
}

// StudentSaveAnswer saves a single answer for a question.
func (h *UserHandler) StudentSaveAnswer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)

		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid exam ID", toast.VariantError)
			return
		}

		questionIDStr := ctx.PostForm("question_id")
		choiceIDStr := ctx.PostForm("choice_id")

		questionID, err := uuid.Parse(questionIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid question", toast.VariantError)
			return
		}

		choiceID, err := uuid.Parse(choiceIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid choice", toast.VariantError)
			return
		}

		attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), examID, userID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "No active attempt found", toast.VariantError)
			return
		}

		if attempt.Status != domain.AttemptStatusInProgress {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Exam already submitted", toast.VariantError)
			return
		}

		// Determine if the choice is correct
		choices, err := h.App.ListChoicesByQuestion(ctx.Request.Context(), questionID)
		isCorrect := false
		if err == nil {
			for _, choice := range choices {
				if choice.ID == choiceID && choice.IsCorrect {
					isCorrect = true
					break
				}
			}
		}

		_, err = h.App.SaveAnswer(ctx.Request.Context(), attempt.ID, questionID, choiceID, isCorrect)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Failed to save answer", toast.VariantError)
			return
		}

		ctx.Header("HX-Reswap", "none")
		ctx.Status(200)
	}
}

// StudentSubmitExam submits the exam attempt and calculates score.
func (h *UserHandler) StudentSubmitExam() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)

		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Invalid exam ID", toast.VariantError)
			return
		}

		attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), examID, userID)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "No active attempt found", toast.VariantError)
			return
		}

		if attempt.Status != domain.AttemptStatusInProgress {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Exam already submitted", toast.VariantError)
			return
		}

		// Calculate score
		answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
		questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), examID)

		// Build a map of question ID to marks
		questionMarks := make(map[uuid.UUID]float64)
		var totalQuestionMarks float64
		for _, q := range questions {
			questionMarks[q.ID] = q.Marks
			totalQuestionMarks += q.Marks
		}

		var rawScore float64
		for _, answer := range answers {
			if answer.IsCorrect != nil && *answer.IsCorrect {
				if marks, ok := questionMarks[answer.QuestionID]; ok {
					rawScore += marks
				}
			}
		}

		exam, _ := h.App.GetExamByID(ctx.Request.Context(), examID)

		var totalScore float64
		if totalQuestionMarks > 0 {
			totalScore = (rawScore / totalQuestionMarks) * exam.TotalMarks
		}

		err = h.App.SubmitExamAttempt(ctx.Request.Context(), attempt.ID, totalScore)
		if err != nil {
			ctx.Header("HX-Reswap", "none")
			helpers.Toast(ctx, "Error", "Failed to submit exam: "+err.Error(), toast.VariantError)
			return
		}

		// Check if this was the last exam for the subject
		totalExams, _ := h.App.CountPublishedExamsBySubject(ctx.Request.Context(), exam.SubjectID)
		submittedExams, _ := h.App.CountSubmittedAttemptsBySubjectForStudent(ctx.Request.Context(), exam.SubjectID, userID)

		if totalExams > 0 && submittedExams >= totalExams {
			ctx.Header("HX-Redirect", fmt.Sprintf("/student/subject/%s/result", exam.SubjectID.String()))
		} else {
			ctx.Header("HX-Redirect", fmt.Sprintf("/student/subject/%s", exam.SubjectID.String()))
		}
	}
}

// StudentAutoSubmit auto-submits exam on timeout or exit (called via sendBeacon).
func (h *UserHandler) StudentAutoSubmit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)

		examIDStr := ctx.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			ctx.Status(400)
			return
		}

		attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), examID, userID)
		if err != nil || attempt.Status != domain.AttemptStatusInProgress {
			ctx.Status(200)
			return
		}

		// Calculate score
		answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
		questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), examID)

		questionMarks := make(map[uuid.UUID]float64)
		var totalQuestionMarks float64
		for _, q := range questions {
			questionMarks[q.ID] = q.Marks
			totalQuestionMarks += q.Marks
		}

		var rawScore float64
		for _, answer := range answers {
			if answer.IsCorrect != nil && *answer.IsCorrect {
				if marks, ok := questionMarks[answer.QuestionID]; ok {
					rawScore += marks
				}
			}
		}

		exam, _ := h.App.GetExamByID(ctx.Request.Context(), examID)

		var totalScore float64
		if totalQuestionMarks > 0 {
			totalScore = (rawScore / totalQuestionMarks) * exam.TotalMarks
		}

		_ = h.App.SubmitExamAttempt(ctx.Request.Context(), attempt.ID, totalScore)
		ctx.Status(200)
	}
}

// StudentSubjectAutoSubmit auto-submits all in-progress exams for a given subject.
func (h *UserHandler) StudentSubjectAutoSubmit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)

		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Status(400)
			return
		}

		exams, _ := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
		for _, exam := range exams {
			if exam.Status != domain.ExamStatusPublished {
				continue
			}
			attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), exam.ID, userID)
			if err != nil || attempt.Status != domain.AttemptStatusInProgress {
				continue
			}

			// Calculate score
			answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
			questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), exam.ID)

			questionMarks := make(map[uuid.UUID]float64)
			var totalQuestionMarks float64
			for _, q := range questions {
				questionMarks[q.ID] = q.Marks
				totalQuestionMarks += q.Marks
			}

			var rawScore float64
			for _, answer := range answers {
				if answer.IsCorrect != nil && *answer.IsCorrect {
					if marks, ok := questionMarks[answer.QuestionID]; ok {
						rawScore += marks
					}
				}
			}

			var totalScore float64
			if totalQuestionMarks > 0 {
				totalScore = (rawScore / totalQuestionMarks) * exam.TotalMarks
			}

			_ = h.App.SubmitExamAttempt(ctx.Request.Context(), attempt.ID, totalScore)
		}

		ctx.Status(200)
	}
}

// StudentSubjectResult renders the final result page for a subject.
func (h *UserHandler) StudentSubjectResult() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, fullname, userID := parseUsername(ctx)
		if username == "" {
			ctx.Redirect(302, "/login")
			return
		}

		studentUser := domain.User{
			ID:       userID,
			Username: username,
			FullName: fullname,
		}

		subjectIDStr := ctx.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		subject, err := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		if err != nil {
			ctx.Redirect(302, "/student/dashboard")
			return
		}

		// Verify all exams are submitted
		totalExams, _ := h.App.CountPublishedExamsBySubject(ctx.Request.Context(), subjectID)
		submittedExams, _ := h.App.CountSubmittedAttemptsBySubjectForStudent(ctx.Request.Context(), subjectID, userID)

		if totalExams == 0 || submittedExams < totalExams {
			ctx.Redirect(302, fmt.Sprintf("/student/subject/%s", subjectID.String()))
			return
		}

		// Get all exams and their attempts for the subject
		exams, _ := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
		var examResults []studentPages.ExamResult
		var totalSubjectScore float64
		var maxPossibleScore float64

		for _, exam := range exams {
			if exam.Status != domain.ExamStatusPublished {
				continue
			}
			attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), exam.ID, userID)
			score := float64(0)
			correctQuestions := 0
			totalQuestions := 0
			timeSpentMinutes := float64(0)

			var questionResults []studentPages.QuestionResult
			answersMap := make(map[uuid.UUID]domain.StudentAnswer)

			if err == nil {
				if attempt.Score != nil {
					score = *attempt.Score
				}
				if attempt.SubmittedAt != nil {
					timeSpentMinutes = attempt.SubmittedAt.Sub(attempt.StartedAt).Minutes()
				}
				answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
				for _, a := range answers {
					answersMap[a.QuestionID] = a
					if a.IsCorrect != nil && *a.IsCorrect {
						correctQuestions++
					}
				}
			}

			questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), exam.ID)
			shuffleQuestionsAndChoices(userID, exam.ID, questions)
			totalQuestions = len(questions)

			for i, q := range questions {
				isCorrect := false
				if ans, ok := answersMap[q.ID]; ok {
					if ans.IsCorrect != nil && *ans.IsCorrect {
						isCorrect = true
					}
				}
				title := q.QuestionTitle
				if title == "" {
					title = fmt.Sprintf("Question %d", i+1)
				} else {
					title = fmt.Sprintf("Question %d: %s", i+1, title)
				}
				questionResults = append(questionResults, studentPages.QuestionResult{
					Title:     title,
					IsCorrect: isCorrect,
				})
			}

			examResults = append(examResults, studentPages.ExamResult{
				Title:            exam.Title,
				Score:            score,
				TotalMarks:       exam.TotalMarks,
				CorrectQuestions: correctQuestions,
				TotalQuestions:   totalQuestions,
				TimeSpentMinutes: timeSpentMinutes,
				Questions:        questionResults,
			})
			totalSubjectScore += score
			maxPossibleScore += exam.TotalMarks
		}

		ctx.Header("Content-Type", "text/html")
		render.Render(ctx, pages.BasePage("Results - "+subject.Title, studentPages.ResultPage(studentPages.ResultParam{
			Student:           studentUser,
			Subject:           subject,
			ExamResults:       examResults,
			TotalSubjectScore: totalSubjectScore,
			MaxPossibleScore:  maxPossibleScore,
		})))
	}
}

// StudentSubjectPDF generates and downloads the result PDF.
func (h *UserHandler) StudentSubjectPDF() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)

		subjectIDStr := ctx.Param("id")
		subjectID, _ := uuid.Parse(subjectIDStr)

		subject, _ := h.App.GetSubjectByID(ctx.Request.Context(), subjectID)
		user, _ := h.App.GetUserByID(ctx.Request.Context(), userID)

		exams, _ := h.App.ListExamsBySubject(ctx.Request.Context(), subjectID)
		var totalScore float64
		var maxScore float64

		m := maroto.New(config.NewBuilder().Build())

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
					text.New("EXAM REPORT", props.Text{
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

		// Student Info Box
		m.AddRows(
			row.New(8).Add(
				col.New(12).Add(
					text.New("STUDENT INFORMATION", props.Text{Size: 7, Style: fontstyle.Bold, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Left: 3, Top: 2}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
			row.New(18).Add(
				col.New(4).Add(
					text.New("Full Name", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Left: 3}),
					text.New(user.FullName, props.Text{Size: 10, Style: fontstyle.Bold, Top: 4, Left: 3}),
				),
				col.New(4).Add(
					text.New("Student ID", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(user.Username, props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
				col.New(4).Add(
					text.New("Report Date", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}}),
					text.New(time.Now().Format("Jan 02, 2006 15:04"), props.Text{Size: 10, Style: fontstyle.Bold, Top: 4}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
		)

		m.AddRows(row.New(15))

		// Exam Table Header
		m.AddRows(
			row.New(10).Add(
				col.New(5).Add(text.New("EXAM COMPONENT", props.Text{Size: 9, Style: fontstyle.Bold, Color: &props.Color{Red: 255, Green: 255, Blue: 255}, Left: 3})),
				col.New(2).Add(text.New("DURATION", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(2).Add(text.New("SCORE", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}})),
				col.New(3).Add(text.New("WEIGHT", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Center, Color: &props.Color{Red: 255, Green: 255, Blue: 255}, Right: 3})),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 30, Green: 41, Blue: 59}}), // Slate 800
		)

		for _, exam := range exams {
			if exam.Status != domain.ExamStatusPublished {
				continue
			}
			attempt, err := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), exam.ID, userID)
			score := float64(0)
			timeSpent := "N/A"

			if err == nil {
				if attempt.Score != nil {
					score = *attempt.Score
				}
				if attempt.SubmittedAt != nil {
					duration := attempt.SubmittedAt.Sub(attempt.StartedAt)
					timeSpent = fmt.Sprintf("%dm %ds", int(duration.Minutes()), int(duration.Seconds())%60)
				}
			}
			totalScore += score
			maxScore += exam.TotalMarks

			// Exam Title Row
			m.AddRows(
				row.New(12).Add(
					col.New(5).Add(text.New(exam.Title, props.Text{Size: 10, Style: fontstyle.Bold, Top: 2, Left: 3})),
					col.New(2).Add(text.New(timeSpent, props.Text{Align: align.Center, Size: 9, Top: 2})),
					col.New(2).Add(text.New(fmt.Sprintf("%.2f", score), props.Text{Align: align.Center, Size: 10, Style: fontstyle.Bold, Top: 2, Color: &props.Color{Red: 30, Green: 58, Blue: 138}})),
					col.New(3).Add(text.New(fmt.Sprintf("%.2f marks", exam.TotalMarks), props.Text{Align: align.Center, Size: 9, Top: 2, Color: &props.Color{Red: 71, Green: 85, Blue: 105}, Right: 3})),
				).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 241, Green: 245, Blue: 249}}), // Slate 100
			)

			// Detailed question breakdown
			if err == nil {
				answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
				questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), exam.ID)
				shuffleQuestionsAndChoices(userID, exam.ID, questions)

				answersMap := make(map[uuid.UUID]domain.StudentAnswer)
				for _, a := range answers {
					answersMap[a.QuestionID] = a
				}

				for i, q := range questions {
					status := "Incorrect"
					qScore := 0.0
					statusColor := &props.Color{Red: 185, Green: 28, Blue: 28} // Red 700
					if ans, ok := answersMap[q.ID]; ok && ans.IsCorrect != nil && *ans.IsCorrect {
						status = "Correct"
						qScore = q.Marks
						statusColor = &props.Color{Red: 21, Green: 128, Blue: 61} // Green 700
					}

					title := q.QuestionTitle
					if title == "" {
						title = fmt.Sprintf("Question %d", i+1)
					} else {
						title = fmt.Sprintf("Question %d: %s", i+1, title)
					}

					rowStyle := &props.Cell{BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255}}
					if i%2 == 1 {
						rowStyle.BackgroundColor = &props.Color{Red: 250, Green: 250, Blue: 250}
					}

					m.AddRows(
						row.New(8).Add(
							col.New(1).Add(text.New("")),
							col.New(7).Add(text.New(title, props.Text{Size: 8, Color: &props.Color{Red: 51, Green: 65, Blue: 85}, Top: 2, Left: 2})),
							col.New(2).Add(text.New(status, props.Text{Size: 8, Align: align.Center, Color: statusColor, Top: 2})),
							col.New(2).Add(text.New(fmt.Sprintf("%.2f / %.2f", qScore, q.Marks), props.Text{Size: 8, Align: align.Center, Color: statusColor, Top: 2, Right: 3})),
						).WithStyle(rowStyle),
					)
				}
			}
			m.AddRows(row.New(5))
		}

		m.AddRows(row.New(10))

		// Summary Section
		m.AddRows(
			row.New(20).Add(
				col.New(7).Add(text.New("FINAL ASSESSMENT SUMMARY", props.Text{Size: 10, Style: fontstyle.Bold, Top: 6, Left: 3})),
				col.New(2).Add(
					text.New("TOTAL SCORE", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Align: align.Right, Top: 3}),
					text.New(fmt.Sprintf("%.2f", totalScore), props.Text{Size: 12, Style: fontstyle.Bold, Align: align.Right, Top: 10, Color: &props.Color{Red: 15, Green: 23, Blue: 42}}),
				),
				col.New(3).Add(
					text.New("MAX POSSIBLE", props.Text{Size: 7, Color: &props.Color{Red: 100, Green: 116, Blue: 139}, Align: align.Right, Top: 3, Right: 3}),
					text.New(fmt.Sprintf("/ %.2f", maxScore), props.Text{Size: 12, Style: fontstyle.Bold, Align: align.Right, Top: 10, Color: &props.Color{Red: 71, Green: 85, Blue: 105}, Right: 3}),
				),
			).WithStyle(&props.Cell{BackgroundColor: &props.Color{Red: 248, Green: 250, Blue: 252}}),
		)

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
			row.New(30).Add(
				col.New(12).Add(
					text.New("VERIFIED BY UOT EXAMINATION MANAGEMENT SYSTEM", props.Text{
						Size:  7,
						Align: align.Center,
						Style: fontstyle.Bold,
						Top:   2,
						Color: &props.Color{Red: 30, Green: 58, Blue: 138},
					}),
				),
			),
		)

		document, _ := m.Generate()

		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_result.pdf", url.PathEscape(subject.Title)))
		ctx.Header("Content-Type", "application/pdf")
		ctx.Data(200, "application/pdf", document.GetBytes())
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// StudentSubjectTrackerWS handles real-time timer sync and auto-submit over WebSocket.
func (h *UserHandler) StudentSubjectTrackerWS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, userID := parseUsername(ctx)
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
			if err != nil {
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

			now := time.Now()
			var remaining int
			if examEndTime != nil {
				remaining = int(examEndTime.Sub(now).Seconds())
			} else {
				remaining = -1 // Force auto-submit if no published exams (e.g. admin ended timer)
			}

			if remaining <= 0 {
				// Time up or manually ended! Force submit
				for _, exam := range exams {
					if exam.Status != domain.ExamStatusPublished && exam.Status != domain.ExamStatusClosed {
						continue
					}
					attempt, attemptErr := h.App.GetAttemptByExamAndStudent(ctx.Request.Context(), exam.ID, userID)
					if attemptErr != nil || attempt.Status != domain.AttemptStatusInProgress {
						continue
					}

					answers, _ := h.App.ListAnswersByAttempt(ctx.Request.Context(), attempt.ID)
					questions, _ := h.App.ListQuestionsByExam(ctx.Request.Context(), exam.ID)

					questionMarks := make(map[uuid.UUID]float64)
					for _, q := range questions {
						questionMarks[q.ID] = q.Marks
					}

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
					// client disconnected
					break
				}
			}
		}
	}
}

func shuffleQuestionsAndChoices(userID, examID uuid.UUID, questions []domain.Question) {
	h := fnv.New64a()
	h.Write(userID[:])
	h.Write(examID[:])
	seed := int64(h.Sum64())

	rng := rand.New(rand.NewSource(seed))

	rng.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	for i := range questions {
		if len(questions[i].Choices) > 0 {
			choices := questions[i].Choices
			rng.Shuffle(len(choices), func(a, b int) {
				choices[a], choices[b] = choices[b], choices[a]
			})
		}
	}
}
