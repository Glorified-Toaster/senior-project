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
SELECT * FROM exams WHERE id = $1;

-- name: ListExamsBySubject :many
SELECT * FROM exams WHERE subject_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: ListAllExams :many
SELECT * FROM exams WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountExams :one
SELECT COUNT(*) FROM exams WHERE deleted_at IS NULL;

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
SELECT * FROM exam_attempts WHERE student_id = $1 ORDER BY started_at DESC;


-- name: SaveAnswer :one
INSERT INTO student_answers (attempt_id, question_id, selected_choice_id, is_correct)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListAnswersByAttempt :many
SELECT * FROM student_answers WHERE attempt_id = $1;

-- name: SearchExams :many
SELECT * FROM exams 
WHERE (title ILIKE '%' || $1::text || '%' OR description ILIKE '%' || $1::text || '%')
AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SoftDeleteExam :exec
UPDATE exams SET deleted_at = NOW() WHERE id = $1;
