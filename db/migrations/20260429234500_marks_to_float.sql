-- +goose Up
-- +goose StatementBegin

-- Alter subjects table
ALTER TABLE subjects ALTER COLUMN total_marks TYPE DOUBLE PRECISION;
ALTER TABLE subjects ALTER COLUMN pass_score TYPE DOUBLE PRECISION;

-- Drop triggers that depend on columns we're about to change
DROP TRIGGER IF EXISTS sync_subject_total_marks_on_update ON exams;
DROP TRIGGER IF EXISTS trigger_redistribute_on_exam_update ON exams;

-- Alter exams table
ALTER TABLE exams ALTER COLUMN total_marks TYPE DOUBLE PRECISION;

-- Recreate triggers for exams
CREATE TRIGGER sync_subject_total_marks_on_update
AFTER UPDATE OF total_marks, deleted_at, subject_id, status ON exams
FOR EACH ROW
EXECUTE FUNCTION sync_subject_total_marks();

CREATE TRIGGER trigger_redistribute_on_exam_update
AFTER UPDATE OF total_marks ON exams
FOR EACH ROW
EXECUTE FUNCTION redistribute_exam_marks();

-- Alter questions table
ALTER TABLE questions ALTER COLUMN marks TYPE DOUBLE PRECISION;

-- Alter exam_attempts table
ALTER TABLE exam_attempts ALTER COLUMN score TYPE DOUBLE PRECISION;

-- Update redistribute_exam_marks function to support float marks
CREATE OR REPLACE FUNCTION redistribute_exam_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_exam_id UUID;
    v_total_marks DOUBLE PRECISION;
    v_q_count INT;
BEGIN
    IF pg_trigger_depth() > 1 THEN
        RETURN NULL;
    END IF;

    IF TG_TABLE_NAME = 'exams' THEN
        target_exam_id := NEW.id;
        v_total_marks := NEW.total_marks;
    ELSIF TG_TABLE_NAME = 'questions' THEN
        IF TG_OP = 'DELETE' THEN
            target_exam_id := OLD.exam_id;
        ELSE
            target_exam_id := NEW.exam_id;
        END IF;
        
        SELECT total_marks INTO v_total_marks FROM exams WHERE id = target_exam_id;
    END IF;

    IF target_exam_id IS NULL OR v_total_marks IS NULL THEN
        RETURN NULL;
    END IF;

    SELECT COUNT(*) INTO v_q_count 
    FROM questions 
    WHERE exam_id = target_exam_id AND deleted_at IS NULL;

    IF v_q_count > 0 THEN
        -- Simply divide total marks by question count
        UPDATE questions 
        SET marks = v_total_marks / v_q_count
        WHERE exam_id = target_exam_id AND deleted_at IS NULL;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert redistribute_exam_marks function
CREATE OR REPLACE FUNCTION redistribute_exam_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_exam_id UUID;
    v_total_marks INT;
    v_q_count INT;
    v_base_mark INT;
    v_remainder INT;
BEGIN
    IF pg_trigger_depth() > 1 THEN
        RETURN NULL;
    END IF;

    IF TG_TABLE_NAME = 'exams' THEN
        target_exam_id := NEW.id;
        v_total_marks := NEW.total_marks;
    ELSIF TG_TABLE_NAME = 'questions' THEN
        IF TG_OP = 'DELETE' THEN
            target_exam_id := OLD.exam_id;
        ELSE
            target_exam_id := NEW.exam_id;
        END IF;
        
        SELECT total_marks INTO v_total_marks FROM exams WHERE id = target_exam_id;
    END IF;

    IF target_exam_id IS NULL OR v_total_marks IS NULL THEN
        RETURN NULL;
    END IF;

    SELECT COUNT(*) INTO v_q_count 
    FROM questions 
    WHERE exam_id = target_exam_id AND deleted_at IS NULL;

    IF v_q_count > 0 THEN
        v_base_mark := v_total_marks / v_q_count;
        v_remainder := v_total_marks % v_q_count;

        WITH ranked_questions AS (
            SELECT id, ROW_NUMBER() OVER (ORDER BY created_at ASC, id ASC) as row_num
            FROM questions
            WHERE exam_id = target_exam_id AND deleted_at IS NULL
        )
        UPDATE questions q
        SET marks = CASE 
            WHEN r.row_num <= v_remainder THEN v_base_mark + 1
            ELSE v_base_mark
        END
        FROM ranked_questions r
        WHERE q.id = r.id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Alter exam_attempts table back to INT
ALTER TABLE exam_attempts ALTER COLUMN score TYPE INT;

-- Alter questions table back to INT
ALTER TABLE questions ALTER COLUMN marks TYPE INT;

-- Alter exams table back to INT
ALTER TABLE exams ALTER COLUMN total_marks TYPE INT;

-- Alter subjects table back to INT
ALTER TABLE subjects ALTER COLUMN total_marks TYPE INT;
ALTER TABLE subjects ALTER COLUMN pass_score TYPE INT;

-- +goose StatementEnd
