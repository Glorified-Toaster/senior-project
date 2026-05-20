-- name: CreateQuestion :one
INSERT INTO questions (exam_id, question_title, question_text, question_type, marks, question_image, checksum)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateQuestion :one
UPDATE questions
SET 
  question_title = $2,
  question_text = $3,
  question_type = $4,
  marks = $5,
  question_image = $6,
  checksum = $7
WHERE id = $1
RETURNING *;

-- name: GetQuestionByID :one
SELECT * FROM questions WHERE deleted_at IS NULL AND id = $1;

-- name: ListQuestionsByExam :many
SELECT * FROM questions WHERE deleted_at IS NULL AND exam_id = $1 ORDER BY created_at ASC;

-- name: CountQuestions :one
SELECT COUNT(*) FROM questions WHERE deleted_at IS NULL;

-- name: SoftDeleteQuestion :exec
UPDATE questions SET deleted_at = NOW() WHERE id = $1;

-- name: RestoreQuestion :exec
UPDATE questions SET deleted_at = NULL WHERE id = $1;

-- name: CreateChoice :one
INSERT INTO choices (question_id, choice_text, is_correct)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateChoices :many
INSERT INTO choices (question_id, choice_text, is_correct)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListChoicesByQuestion :many
SELECT * FROM choices WHERE deleted_at IS NULL AND question_id = $1 ORDER BY created_at ASC;

-- name: ListChoicesByQuestionSeeded :many
SELECT * FROM choices 
WHERE deleted_at IS NULL 
  AND question_id = $1 
ORDER BY md5(id::text || $2::text);

-- name: GetQuestionByChecksum :one
SELECT * FROM questions WHERE deleted_at IS NULL AND checksum = $1;

-- name: DeleteChoicesByQuestion :exec
UPDATE choices SET deleted_at = NOW() WHERE question_id = $1;

-- name: DeleteQuestionByID :exec
UPDATE questions SET deleted_at = NOW() WHERE id = $1;

-- name: SetQuestionImagePath :exec
UPDATE questions SET question_image = $2 WHERE deleted_at IS NULL AND question_type = 'IMAGE' AND id = $1;

-- name: UpdateChoice :one
UPDATE choices
SET choice_text = $2, is_correct = $3, deleted_at = NULL, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListAllChoicesByQuestion :many
SELECT * FROM choices WHERE question_id = $1 ORDER BY created_at ASC;

-- name: SoftDeleteChoiceByID :exec
UPDATE choices SET deleted_at = NOW() WHERE id = $1;