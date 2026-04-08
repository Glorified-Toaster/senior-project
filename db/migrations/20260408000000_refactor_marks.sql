-- +goose Up
-- +goose StatementBegin

-- =====================================================================
-- Remove pass_score from exams (cumulative grading lives on subjects)
-- =====================================================================
ALTER TABLE exams DROP COLUMN IF EXISTS pass_score;

-- =====================================================================
-- Relax constraints on subjects so total_marks can be 0 (no exams yet)
-- =====================================================================
ALTER TABLE subjects DROP CONSTRAINT IF EXISTS subjects_total_marks_check;
ALTER TABLE subjects DROP CONSTRAINT IF EXISTS subjects_check;

-- total_marks is now auto-computed; allow 0 as a starting value
ALTER TABLE subjects ALTER COLUMN total_marks SET DEFAULT 0;
ALTER TABLE subjects ADD CONSTRAINT subjects_total_marks_check CHECK (total_marks >= 0);

-- pass_score must be >= 0 (0 means "no threshold set yet")
ALTER TABLE subjects ADD CONSTRAINT subjects_pass_score_check CHECK (pass_score >= 0);

-- =====================================================================
-- Trigger: keep subjects.total_marks in sync with SUM(exams.total_marks)
-- =====================================================================
CREATE OR REPLACE FUNCTION sync_subject_total_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_subject_id UUID;
BEGIN
    -- Determine which subject to update
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

CREATE TRIGGER sync_subject_total_marks
AFTER INSERT OR DELETE ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

-- Separate trigger for UPDATE to only fire when relevant columns change
CREATE TRIGGER sync_subject_total_marks_on_update
AFTER UPDATE OF total_marks, deleted_at, subject_id ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

-- =====================================================================
-- Backfill: recompute total_marks for all existing subjects
-- =====================================================================
UPDATE subjects s
SET total_marks = COALESCE((
    SELECT SUM(e.total_marks)
    FROM exams e
    WHERE e.subject_id = s.id
      AND e.deleted_at IS NULL
), 0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS sync_subject_total_marks_on_update ON exams;
DROP TRIGGER IF EXISTS sync_subject_total_marks ON exams;
DROP FUNCTION IF EXISTS sync_subject_total_marks();

ALTER TABLE subjects DROP CONSTRAINT IF EXISTS subjects_total_marks_check;
ALTER TABLE subjects DROP CONSTRAINT IF EXISTS subjects_pass_score_check;
ALTER TABLE subjects ADD CONSTRAINT subjects_total_marks_check CHECK (total_marks > 0);
ALTER TABLE subjects ADD CONSTRAINT subjects_pass_score_check CHECK (pass_score > 0 AND pass_score <= total_marks);

ALTER TABLE exams ADD COLUMN pass_score INT NOT NULL DEFAULT 0 CHECK (pass_score >= 0);

-- +goose StatementEnd
