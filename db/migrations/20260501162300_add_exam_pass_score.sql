-- +goose Up
ALTER TABLE exams ADD COLUMN pass_score FLOAT NOT NULL DEFAULT 0;

-- Set default pass_score to total_marks / 2 for existing exams
UPDATE exams SET pass_score = total_marks / 2.0;

-- +goose Down
ALTER TABLE exams DROP COLUMN pass_score;
