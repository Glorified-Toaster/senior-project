-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION sync_subject_total_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_subject_id UUID;
BEGIN
    IF TG_OP = 'DELETE' THEN
        target_subject_id := OLD.subject_id;
    ELSE
        target_subject_id := NEW.subject_id;
    END IF;

    UPDATE subjects
    SET total_marks = COALESCE((
        SELECT SUM(total_marks)
        FROM exams
        WHERE subject_id = target_subject_id
          AND deleted_at IS NULL
    ), 0)
    WHERE id = target_subject_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger an update on all subjects to recalculate their total_marks
UPDATE subjects
SET total_marks = COALESCE((
    SELECT SUM(total_marks)
    FROM exams
    WHERE subject_id = subjects.id
      AND deleted_at IS NULL
), 0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION sync_subject_total_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_subject_id UUID;
BEGIN
    IF TG_OP = 'DELETE' THEN
        target_subject_id := OLD.subject_id;
    ELSE
        target_subject_id := NEW.subject_id;
    END IF;

    UPDATE subjects
    SET total_marks = COALESCE((
        SELECT SUM(total_marks)
        FROM exams
        WHERE subject_id = target_subject_id
          AND deleted_at IS NULL
          AND status = 'PUBLISHED'
    ), 0)
    WHERE id = target_subject_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

UPDATE subjects
SET total_marks = COALESCE((
    SELECT SUM(total_marks)
    FROM exams
    WHERE subject_id = subjects.id
      AND deleted_at IS NULL
      AND status = 'PUBLISHED'
), 0);

-- +goose StatementEnd
