package repository

import (
	"context"
	"uot-exam/internal/adapters/outbound/database"
	"uot-exam/internal/adapters/outbound/database/sqlc"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ExamRepository struct {
	queries *sqlc.Queries
}

func NewExamRepository(queries *sqlc.Queries) *ExamRepository {
	return &ExamRepository{queries: queries}
}

func (r *ExamRepository) ListAll(ctx context.Context, arg ports.ListAllExamsParams) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	claimedExams, err := queries.ListAllExams(ctx, sqlc.ListAllExamsParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range claimedExams {
		exams = append(exams, mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions))
	}

	return exams, nil
}

func (r *ExamRepository) Search(ctx context.Context, arg ports.SearchExamsParams) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	claimedExams, err := queries.SearchExams(ctx, sqlc.SearchExamsParams{
		Column1: arg.Search,
		Limit:   arg.Limit,
		Offset:  arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range claimedExams {
		exams = append(exams, mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions))
	}

	return exams, nil
}

func (r *ExamRepository) Count(ctx context.Context) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountExams(ctx)
}

func (r *ExamRepository) CountSearch(ctx context.Context, search string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountSearchExams(ctx, search)
}

func (r *ExamRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.SoftDeleteExam(ctx, id)
}

func mapSqlcExamToDomain(exam sqlc.Exam) domain.Exam {
	return domain.Exam{
		ID:          exam.ID,
		SubjectID:   exam.SubjectID.UUID,
		Title:       exam.Title,
		Description: exam.Description,
		TotalMarks:  exam.TotalMarks,
		PassScore:   exam.PassScore,
		Status:      domain.ExamStatus(exam.Status),
		CreatedBy:   exam.CreatedBy.UUID,
		CreatedAt:   exam.CreatedAt.Time,
		UpdatedAt:   exam.UpdatedAt.Time,
		DeletedAt:   toTimePtr(exam.DeletedAt),
	}
}

func mapSqlcFullExamToDomain(exam sqlc.Exam, totalQuestions int64) domain.Exam {
	d := mapSqlcExamToDomain(exam)
	d.TotalQuestions = totalQuestions
	return d
}

func (r *ExamRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exam, err := queries.GetExamByID(ctx, id)
	if err != nil {
		return domain.Exam{}, err
	}

	return mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions), nil
}

func (r *ExamRepository) Create(ctx context.Context, arg ports.CreateExamParams) (domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exam, err := queries.CreateExam(ctx, sqlc.CreateExamParams{
		SubjectID:   uuid.NullUUID{UUID: arg.SubjectID, Valid: true},
		Title:       arg.Title,
		Description: arg.Description,
		TotalMarks:  arg.TotalMarks,
		PassScore:   arg.PassScore,
		Status:      sqlc.ExamStatusType(arg.Status),
		CreatedBy:   uuid.NullUUID{UUID: arg.CreatedBy, Valid: true},
	})
	if err != nil {
		return domain.Exam{}, err
	}

	return mapSqlcExamToDomain(exam), nil
}

func (r *ExamRepository) ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.ListExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions))
	}

	return exams, nil
}

func (r *ExamRepository) Update(ctx context.Context, arg ports.UpdateExamParams) (domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exam, err := queries.UpdateExam(ctx, sqlc.UpdateExamParams{
		ID:          arg.ID,
		Title:       arg.Title,
		Description: arg.Description,
		TotalMarks:  arg.TotalMarks,
		PassScore:   arg.PassScore,
		Status:      sqlc.ExamStatusType(arg.Status),
	})
	if err != nil {
		return domain.Exam{}, err
	}

	return mapSqlcExamToDomain(exam), nil
}

func (r *ExamRepository) SearchBySubject(ctx context.Context, arg ports.SearchExamsBySubjectParams) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.SearchExamsBySubject(ctx, sqlc.SearchExamsBySubjectParams{
		SubjectID: uuid.NullUUID{UUID: arg.SubjectID, Valid: true},
		Column2:   arg.Search,
		Limit:     arg.Limit,
		Offset:    arg.Offset,
	})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions))
	}

	return exams, nil
}

func (r *ExamRepository) CountSearchBySubject(ctx context.Context, subjectID uuid.UUID, search string) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountSearchExamsBySubject(ctx, sqlc.CountSearchExamsBySubjectParams{
		SubjectID: uuid.NullUUID{UUID: subjectID, Valid: true},
		Column2:   search,
	})
}

func (r *ExamRepository) ListExamsForStudent(ctx context.Context, studentID uuid.UUID) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.ListExamsForStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions))
	}

	return exams, nil
}

func (r *ExamRepository) ListExamsCreatedBy(ctx context.Context, instructorID uuid.UUID) ([]domain.Exam, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcExams, err := queries.ListExamsCreatedBy(ctx, uuid.NullUUID{UUID: instructorID, Valid: true})
	if err != nil {
		return nil, err
	}

	var exams []domain.Exam
	for _, exam := range sqlcExams {
		exams = append(exams, mapSqlcFullExamToDomain(exam.Exam, exam.TotalQuestions))
	}

	return exams, nil
}

func toFloat64Ptr(val pgtype.Float8) *float64 {
	if !val.Valid {
		return nil
	}
	return &val.Float64
}

func mapSqlcAttemptToDomain(attempt sqlc.ExamAttempt) domain.ExamAttempt {
	return domain.ExamAttempt{
		ID:          attempt.ID,
		ExamID:      attempt.ExamID.UUID,
		StudentID:   attempt.StudentID.UUID,
		StartedAt:   attempt.StartedAt.Time,
		SubmittedAt: toTimePtr(attempt.SubmittedAt),
		Score:       toFloat64Ptr(attempt.Score),
		Status:      domain.AttemptStatus(attempt.Status),
		CreatedAt:   attempt.CreatedAt.Time,
	}
}

func mapSqlcAttemptRowToDomain(attempt sqlc.ListAttemptsByStudentRow) domain.ExamAttempt {
	return domain.ExamAttempt{
		ID:           attempt.ID,
		ExamID:       attempt.ExamID.UUID,
		StudentID:    attempt.StudentID.UUID,
		StartedAt:    attempt.StartedAt.Time,
		SubmittedAt:  toTimePtr(attempt.SubmittedAt),
		Score:        toFloat64Ptr(attempt.Score),
		Status:       domain.AttemptStatus(attempt.Status),
		CreatedAt:    attempt.CreatedAt.Time,
		ExamTitle:    attempt.ExamTitle,
		SubjectTitle: attempt.SubjectTitle,
	}
}

func (r *ExamRepository) ListAttemptsByStudent(ctx context.Context, studentID uuid.UUID) ([]domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcAttempts, err := queries.ListAttemptsByStudent(ctx, uuid.NullUUID{UUID: studentID, Valid: true})
	if err != nil {
		return nil, err
	}

	var attempts []domain.ExamAttempt
	for _, attempt := range sqlcAttempts {
		attempts = append(attempts, mapSqlcAttemptRowToDomain(attempt))
	}

	return attempts, nil
}

func (r *ExamRepository) GetAttemptByExamAndStudent(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	attempt, err := queries.GetAttemptByExamAndStudent(ctx, sqlc.GetAttemptByExamAndStudentParams{
		ExamID:    uuid.NullUUID{UUID: examID, Valid: true},
		StudentID: uuid.NullUUID{UUID: studentID, Valid: true},
	})
	if err != nil {
		return domain.ExamAttempt{}, err
	}

	return mapSqlcAttemptToDomain(attempt), nil
}

func (r *ExamRepository) StartExamAttempt(ctx context.Context, examID uuid.UUID, studentID uuid.UUID) (domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	attempt, err := queries.StartExamAttempt(ctx, sqlc.StartExamAttemptParams{
		ExamID:    uuid.NullUUID{UUID: examID, Valid: true},
		StudentID: uuid.NullUUID{UUID: studentID, Valid: true},
	})
	if err != nil {
		return domain.ExamAttempt{}, err
	}

	return mapSqlcAttemptToDomain(attempt), nil
}

func (r *ExamRepository) SubmitExamAttempt(ctx context.Context, attemptID uuid.UUID, score float64) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.SubmitExamAttempt(ctx, sqlc.SubmitExamAttemptParams{
		ID:    attemptID,
		Score: pgtype.Float8{Float64: score, Valid: true},
	})
}

func mapSqlcAnswerToDomain(answer sqlc.StudentAnswer) domain.StudentAnswer {
	var isCorrect *bool
	if answer.IsCorrect.Valid {
		isCorrect = &answer.IsCorrect.Bool
	}
	return domain.StudentAnswer{
		ID:               answer.ID,
		AttemptID:        answer.AttemptID.UUID,
		QuestionID:       answer.QuestionID.UUID,
		SelectedChoiceID: answer.SelectedChoiceID.UUID,
		IsCorrect:        isCorrect,
		AnsweredAt:       answer.AnsweredAt.Time,
	}
}

func (r *ExamRepository) SaveAnswer(ctx context.Context, attemptID uuid.UUID, questionID uuid.UUID, choiceID uuid.UUID, isCorrect bool) (domain.StudentAnswer, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	answer, err := queries.SaveAnswer(ctx, sqlc.SaveAnswerParams{
		AttemptID:        uuid.NullUUID{UUID: attemptID, Valid: true},
		QuestionID:       uuid.NullUUID{UUID: questionID, Valid: true},
		SelectedChoiceID: uuid.NullUUID{UUID: choiceID, Valid: true},
		IsCorrect:        pgtype.Bool{Bool: isCorrect, Valid: true},
	})
	if err != nil {
		return domain.StudentAnswer{}, err
	}

	return mapSqlcAnswerToDomain(answer), nil
}

func (r *ExamRepository) ListAnswersByAttempt(ctx context.Context, attemptID uuid.UUID) ([]domain.StudentAnswer, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcAnswers, err := queries.ListAnswersByAttempt(ctx, uuid.NullUUID{UUID: attemptID, Valid: true})
	if err != nil {
		return nil, err
	}

	var answers []domain.StudentAnswer
	for _, answer := range sqlcAnswers {
		answers = append(answers, mapSqlcAnswerToDomain(answer))
	}

	return answers, nil
}

func (r *ExamRepository) CountExamsBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}

func (r *ExamRepository) CountPublishedBySubject(ctx context.Context, subjectID uuid.UUID) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountPublishedExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}

func (r *ExamRepository) CountSubmittedAttemptsBySubjectForStudent(ctx context.Context, subjectID uuid.UUID, studentID uuid.UUID) (int64, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.CountSubmittedAttemptsBySubjectForStudent(ctx, sqlc.CountSubmittedAttemptsBySubjectForStudentParams{
		SubjectID: uuid.NullUUID{UUID: subjectID, Valid: true},
		StudentID: uuid.NullUUID{UUID: studentID, Valid: true},
	})
}

func (r *ExamRepository) PublishDraftExamsBySubject(ctx context.Context, subjectID uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.PublishDraftExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}

func (r *ExamRepository) ClosePublishedExamsBySubject(ctx context.Context, subjectID uuid.UUID) error {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	return queries.ClosePublishedExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
}

func (r *ExamRepository) ListInProgressAttemptsByExam(ctx context.Context, examID uuid.UUID) ([]domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcAttempts, err := queries.ListInProgressAttemptsByExam(ctx, uuid.NullUUID{UUID: examID, Valid: true})
	if err != nil {
		return nil, err
	}

	var attempts []domain.ExamAttempt
	for _, attempt := range sqlcAttempts {
		attempts = append(attempts, mapSqlcAttemptToDomain(attempt))
	}

	return attempts, nil
}

func (r *ExamRepository) GetExamAnalytics(ctx context.Context, examID uuid.UUID) (domain.ExamAnalytics, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exam, err := queries.GetExamByID(ctx, examID)
	if err != nil {
		return domain.ExamAnalytics{}, err
	}

	dbAttempts, err := queries.ListAttemptsByExam(ctx, uuid.NullUUID{UUID: examID, Valid: true})
	if err != nil {
		return domain.ExamAnalytics{}, err
	}

	qStats, err := queries.GetExamQuestionAnalytics(ctx, uuid.NullUUID{UUID: examID, Valid: true})
	if err != nil {
		return domain.ExamAnalytics{}, err
	}

	var totalScore float64
	var maxScore float64 = -1.0
	var minScore float64 = -1.0
	var passCount int64
	var failCount int64
	var submittedCount int64

	for _, att := range dbAttempts {
		if att.Status != "SUBMITTED" && att.Status != "GRADED" {
			continue
		}
		submittedCount++
		score := 0.0
		if att.Score.Valid {
			score = att.Score.Float64
		}

		totalScore += score
		if maxScore == -1.0 || score > maxScore {
			maxScore = score
		}
		if minScore == -1.0 || score < minScore {
			minScore = score
		}

		if score >= exam.Exam.PassScore {
			passCount++
		} else {
			failCount++
		}
	}

	if minScore == -1.0 {
		minScore = 0
	}
	if maxScore == -1.0 {
		maxScore = 0
	}

	var avgScore float64
	if submittedCount > 0 {
		avgScore = totalScore / float64(submittedCount)
	}

	var passRate float64
	if submittedCount > 0 {
		passRate = (float64(passCount) / float64(submittedCount)) * 100
	}

	var questionStats []domain.QuestionAnalytics
	for _, qs := range qStats {
		accuracy := 0.0
		if qs.TotalAnswers > 0 {
			accuracy = (float64(qs.CorrectAnswers) / float64(qs.TotalAnswers)) * 100
		}
		questionStats = append(questionStats, domain.QuestionAnalytics{
			QuestionID:     qs.QuestionID,
			QuestionTitle:  qs.QuestionTitle,
			QuestionType:   string(qs.QuestionType),
			MaxMarks:       qs.MaxMarks,
			TotalAnswers:   qs.TotalAnswers,
			CorrectAnswers: qs.CorrectAnswers,
			Accuracy:       accuracy,
		})
	}

	return domain.ExamAnalytics{
		ExamID:        examID,
		TotalAttempts: int64(len(dbAttempts)),
		AverageScore:  avgScore,
		MaxScore:      maxScore,
		MinScore:      minScore,
		PassCount:     passCount,
		FailCount:     failCount,
		PassRate:      passRate,
		QuestionStats: questionStats,
	}, nil
}

func mapSqlcAttemptWithStudentToDomain(attempt sqlc.ListAttemptsByExamRow) domain.ExamAttempt {
	return domain.ExamAttempt{
		ID:              attempt.ID,
		ExamID:          attempt.ExamID.UUID,
		StudentID:       attempt.StudentID.UUID,
		StartedAt:       attempt.StartedAt.Time,
		SubmittedAt:     toTimePtr(attempt.SubmittedAt),
		Score:           toFloat64Ptr(attempt.Score),
		Status:          domain.AttemptStatus(attempt.Status),
		CreatedAt:       attempt.CreatedAt.Time,
		StudentName:     attempt.StudentName,
		StudentUsername: attempt.StudentUsername,
	}
}

func (r *ExamRepository) ListAttemptsByExam(ctx context.Context, examID uuid.UUID) ([]domain.ExamAttempt, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	sqlcAttempts, err := queries.ListAttemptsByExam(ctx, uuid.NullUUID{UUID: examID, Valid: true})
	if err != nil {
		return nil, err
	}

	var attempts []domain.ExamAttempt
	for _, attempt := range sqlcAttempts {
		attempts = append(attempts, mapSqlcAttemptWithStudentToDomain(attempt))
	}
	return attempts, nil
}

func (r *ExamRepository) GetSubjectAnalytics(ctx context.Context, subjectID uuid.UUID) (domain.SubjectAnalytics, error) {
	queries := r.queries
	if tx := database.ExtractTx(ctx); tx != nil {
		queries = queries.WithTx(tx)
	}

	exams, err := queries.ListExamsBySubject(ctx, uuid.NullUUID{UUID: subjectID, Valid: true})
	if err != nil {
		return domain.SubjectAnalytics{}, err
	}

	var totalExams int64 = int64(len(exams))
	var totalAttempts int64
	var totalPassCount int64
	var totalFailCount int64
	var totalScore float64
	var submittedCount int64
	var maxScore float64 = -1.0
	var minScore float64 = -1.0
	var examStats []domain.ExamSummaryAnalytics

	for _, exam := range exams {
		attempts, err := queries.ListAttemptsByExam(ctx, uuid.NullUUID{UUID: exam.Exam.ID, Valid: true})
		if err != nil {
			continue
		}

		totalAttempts += int64(len(attempts))
		var examPassCount int64
		var examSubmittedCount int64

		for _, att := range attempts {
			if att.Status != "SUBMITTED" && att.Status != "GRADED" {
				continue
			}
			examSubmittedCount++
			submittedCount++
			score := 0.0
			if att.Score.Valid {
				score = att.Score.Float64
			}

			totalScore += score
			if maxScore == -1.0 || score > maxScore {
				maxScore = score
			}
			if minScore == -1.0 || score < minScore {
				minScore = score
			}

			if score >= exam.Exam.PassScore {
				examPassCount++
				totalPassCount++
			} else {
				totalFailCount++
			}
		}

		var examPassRate float64
		if examSubmittedCount > 0 {
			examPassRate = (float64(examPassCount) / float64(examSubmittedCount)) * 100
		}

		examStats = append(examStats, domain.ExamSummaryAnalytics{
			ExamID:        exam.Exam.ID,
			Title:         exam.Exam.Title,
			PassRate:      examPassRate,
			TotalAttempts: int64(len(attempts)),
		})
	}

	if minScore == -1.0 {
		minScore = 0
	}
	if maxScore == -1.0 {
		maxScore = 0
	}

	var avgScore float64
	if submittedCount > 0 {
		avgScore = totalScore / float64(submittedCount)
	}

	var passRate float64
	if submittedCount > 0 {
		passRate = (float64(totalPassCount) / float64(submittedCount)) * 100
	}

	return domain.SubjectAnalytics{
		SubjectID:     subjectID,
		TotalExams:    totalExams,
		TotalAttempts: totalAttempts,
		AverageScore:  avgScore,
		MaxScore:      maxScore,
		MinScore:      minScore,
		PassCount:     totalPassCount,
		FailCount:     totalFailCount,
		PassRate:      passRate,
		ExamStats:     examStats,
	}, nil
}
