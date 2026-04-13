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
          AND status = 'PUBLISHED'
    ), 0)
    WHERE id = target_subject_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Update the existing trigger to also fire on status changes
DROP TRIGGER IF EXISTS sync_subject_total_marks_on_update ON exams;
CREATE TRIGGER sync_subject_total_marks_on_update
AFTER UPDATE OF total_marks, deleted_at, subject_id, status ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

-- Recalculate total marks for all existing subjects
UPDATE subjects s
SET total_marks = COALESCE((
    SELECT SUM(e.total_marks)
    FROM exams e
    WHERE e.subject_id = s.id
      AND e.deleted_at IS NULL
      AND e.status = 'PUBLISHED'
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
    ), 0)
    WHERE id = target_subject_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS sync_subject_total_marks_on_update ON exams;
CREATE TRIGGER sync_subject_total_marks_on_update
AFTER UPDATE OF total_marks, deleted_at, subject_id ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

-- Recalculate back to original logic (all non-deleted exams)
UPDATE subjects s
SET total_marks = COALESCE((
    SELECT SUM(e.total_marks)
    FROM exams e
    WHERE e.subject_id = s.id
      AND e.deleted_at IS NULL
), 0);
-- +goose StatementEnd
