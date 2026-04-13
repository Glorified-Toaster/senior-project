-- name: CreateExam :one
INSERT INTO exams (subject_id, title, description, total_marks, status, created_by)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: UpdateExam :one
UPDATE exams
SET 
  title = $2,
  description = $3,
  total_marks = $4,
  status = $5
WHERE id = $1
RETURNING *;

-- name: GetExamByID :one
SELECT * FROM exams WHERE id = $1;

-- name: ListExamsBySubject :many
SELECT * FROM exams WHERE subject_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC;

-- name: ListAllExams :many
SELECT * FROM exams WHERE deleted_at IS NULL ORDER BY created_at DESC, id ASC LIMIT $1 OFFSET $2;

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
ORDER BY created_at DESC, id ASC
LIMIT $2 OFFSET $3;

-- name: SearchExamsBySubject :many
SELECT * FROM exams 
WHERE subject_id = $1
AND (title ILIKE '%' || $2::text || '%' OR description ILIKE '%' || $2::text || '%')
AND deleted_at IS NULL
ORDER BY created_at DESC, id ASC
LIMIT $3 OFFSET $4;

-- name: CountSearchExamsBySubject :one
SELECT COUNT(*) FROM exams 
WHERE subject_id = $1
AND (title ILIKE '%' || $2::text || '%' OR description ILIKE '%' || $2::text || '%')
AND deleted_at IS NULL;

-- name: SoftDeleteExam :exec
UPDATE exams SET deleted_at = NOW() WHERE id = $1;

-- name: ListExamsCreatedBy :many
SELECT * FROM exams 
WHERE created_by = $1 AND deleted_at IS NULL 
ORDER BY created_at DESC;

-- name: ListExamsForStudent :many
SELECT e.* 
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
