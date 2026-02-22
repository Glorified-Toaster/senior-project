-- name: CreateExam :one
INSERT INTO exams (subject_id, title, description, duration_minutes, total_marks, start_time, end_time, status, created_by)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetExamByID :one
SELECT * FROM exams WHERE id = $1;

-- name: ListExamsBySubject :many
SELECT * FROM exams WHERE subject_id = $1 ORDER BY created_at DESC;


-- name: CreateQuestion :one
INSERT INTO questions (exam_id, question_text, marks, position)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListQuestionsByExam :many
SELECT * FROM questions WHERE exam_id = $1 ORDER BY position ASC;


-- name: CreateChoice :one
INSERT INTO choices (question_id, choice_text, is_correct)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListChoicesByQuestion :many
SELECT * FROM choices WHERE question_id = $1;


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
