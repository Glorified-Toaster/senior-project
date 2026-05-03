-- name: CreateExam :one
INSERT INTO exams (subject_id, title, description, total_marks, pass_score, status, created_by)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING *;

-- name: UpdateExam :one
UPDATE exams
SET 
  title = $2,
  description = $3,
  total_marks = $4,
  pass_score = $5,
  status = $6
WHERE id = $1
RETURNING *;

-- name: GetExamByID :one
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e 
WHERE e.id = $1;

-- name: ListExamsBySubject :many
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e 
WHERE e.subject_id = $1 AND e.deleted_at IS NULL 
ORDER BY e.created_at DESC;

-- name: ListAllExams :many
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e 
WHERE e.deleted_at IS NULL 
ORDER BY e.created_at DESC, e.id ASC 
LIMIT $1 OFFSET $2;

-- name: CountExams :one
SELECT COUNT(*) FROM exams WHERE deleted_at IS NULL;

-- name: CountSearchExams :one
SELECT COUNT(*) FROM exams 
WHERE (title ILIKE '%' || $1::text || '%' OR description ILIKE '%' || $1::text || '%')
AND deleted_at IS NULL;

-- name: StartExamAttempt :one
INSERT INTO exam_attempts (exam_id, student_id)
VALUES ($1, $2)
RETURNING *;

-- name: SubmitExamAttempt :exec
UPDATE exam_attempts
SET submitted_at = NOW(), status = 'SUBMITTED', score = $1
WHERE id = $2;

-- name: GetAttemptByID :one
SELECT * FROM exam_attempts WHERE id = $1;

-- name: ListAttemptsByStudent :many
SELECT ea.*, e.title as exam_title, s.title as subject_title
FROM exam_attempts ea
JOIN exams e ON ea.exam_id = e.id
JOIN subjects s ON e.subject_id = s.id
WHERE ea.student_id = $1 
ORDER BY ea.started_at DESC;


-- name: SaveAnswer :one
INSERT INTO student_answers (attempt_id, question_id, selected_choice_id, is_correct)
VALUES ($1, $2, $3, $4)
ON CONFLICT (attempt_id, question_id) 
DO UPDATE SET selected_choice_id = EXCLUDED.selected_choice_id, is_correct = EXCLUDED.is_correct, answered_at = NOW()
RETURNING *;

-- name: ListAnswersByAttempt :many
SELECT * FROM student_answers WHERE attempt_id = $1;

-- name: SearchExams :many
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e 
WHERE (e.title ILIKE '%' || $1::text || '%' OR e.description ILIKE '%' || $1::text || '%')
AND e.deleted_at IS NULL
ORDER BY e.created_at DESC, e.id ASC
LIMIT $2 OFFSET $3;

-- name: SearchExamsBySubject :many
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e 
WHERE e.subject_id = $1
AND (e.title ILIKE '%' || $2::text || '%' OR e.description ILIKE '%' || $2::text || '%')
AND e.deleted_at IS NULL
ORDER BY e.created_at DESC, e.id ASC
LIMIT $3 OFFSET $4;

-- name: CountSearchExamsBySubject :one
SELECT COUNT(*) FROM exams 
WHERE subject_id = $1
AND (title ILIKE '%' || $2::text || '%' OR description ILIKE '%' || $2::text || '%')
AND deleted_at IS NULL;

-- name: SoftDeleteExam :exec
UPDATE exams SET deleted_at = NOW() WHERE id = $1;

-- name: ListExamsCreatedBy :many
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e 
WHERE e.created_by = $1 AND e.deleted_at IS NULL 
ORDER BY e.created_at DESC;

-- name: ListExamsForStudent :many
SELECT 
    sqlc.embed(e),
    (SELECT COUNT(*) FROM questions q WHERE q.exam_id = e.id AND q.deleted_at IS NULL) as total_questions
FROM exams e
JOIN subject_students ss ON e.subject_id = ss.subject_id
WHERE ss.student_id = $1 AND e.deleted_at IS NULL AND ss.deleted_at IS NULL
ORDER BY e.created_at DESC;

-- name: GetAttemptByExamAndStudent :one
SELECT * FROM exam_attempts
WHERE exam_id = $1 AND student_id = $2
ORDER BY started_at DESC
LIMIT 1;

-- name: CountExamsBySubject :one
SELECT COUNT(*) FROM exams
WHERE subject_id = $1 AND deleted_at IS NULL;

-- name: CountSubmittedAttemptsBySubjectForStudent :one
SELECT COUNT(DISTINCT ea.exam_id)
FROM exam_attempts ea
JOIN exams e ON ea.exam_id = e.id
WHERE e.subject_id = $1
  AND ea.student_id = $2
  AND ea.status IN ('SUBMITTED', 'GRADED')
  AND e.deleted_at IS NULL;

-- name: CountPublishedExamsBySubject :one
SELECT COUNT(*) FROM exams
WHERE subject_id = $1 AND deleted_at IS NULL AND status = 'PUBLISHED';

-- name: PublishDraftExamsBySubject :exec
UPDATE exams
SET status = 'PUBLISHED', updated_at = NOW()
WHERE subject_id = $1 AND status = 'DRAFT' AND deleted_at IS NULL;

-- name: ClosePublishedExamsBySubject :exec
UPDATE exams
SET status = 'CLOSED', updated_at = NOW()
WHERE subject_id = $1 AND status = 'PUBLISHED' AND deleted_at IS NULL;

-- name: ListInProgressAttemptsByExam :many
SELECT * FROM exam_attempts
WHERE exam_id = $1 AND status = 'IN_PROGRESS'
ORDER BY started_at DESC;

-- name: ListAttemptsByExam :many
SELECT ea.*, u.full_name as student_name, u.username as student_username
FROM exam_attempts ea
JOIN users u ON ea.student_id = u.id
WHERE ea.exam_id = $1
ORDER BY ea.started_at DESC;

-- name: GetExamQuestionAnalytics :many
SELECT 
    q.id as question_id,
    q.question_title,
    q.question_type,
    q.marks as max_marks,
    COUNT(sa.id) as total_answers,
    COALESCE(SUM(CASE WHEN sa.is_correct = TRUE THEN 1 ELSE 0 END), 0)::bigint as correct_answers
FROM questions q
LEFT JOIN student_answers sa ON q.id = sa.question_id
LEFT JOIN exam_attempts ea ON sa.attempt_id = ea.id AND ea.status IN ('SUBMITTED', 'GRADED')
WHERE q.exam_id = $1 AND q.deleted_at IS NULL
GROUP BY q.id, q.question_title, q.question_type, q.marks, q.created_at
ORDER BY q.created_at ASC;
