-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION redistribute_exam_marks()
RETURNS TRIGGER AS $$
DECLARE
    target_exam_id UUID;
    v_total_marks INT;
    v_q_count INT;
    v_base_mark INT;
    v_remainder INT;
BEGIN
    -- Prevent infinite recursion
    -- This check ensures that the trigger doesn't fire again if it's already 
    -- running for this session. PostgreSQL 10+ handles this well, but depth 
    -- check is a robust way to avoid unintended cascades.
    IF pg_trigger_depth() > 1 THEN
        RETURN NULL;
    END IF;

    -- Determine the exam ID to process
    IF TG_TABLE_NAME = 'exams' THEN
        target_exam_id := NEW.id;
        v_total_marks := NEW.total_marks;
    ELSIF TG_TABLE_NAME = 'questions' THEN
        IF TG_OP = 'DELETE' THEN
            target_exam_id := OLD.exam_id;
        ELSE
            target_exam_id := NEW.exam_id;
        END IF;
        
        -- Get total marks from the exam table
        SELECT total_marks INTO v_total_marks FROM exams WHERE id = target_exam_id;
    END IF;

    -- If no exam found (shouldn't happen with foreign keys but safe check)
    IF target_exam_id IS NULL OR v_total_marks IS NULL THEN
        RETURN NULL;
    END IF;

    -- Count existing non-deleted questions
    SELECT COUNT(*) INTO v_q_count 
    FROM questions 
    WHERE exam_id = target_exam_id AND deleted_at IS NULL;

    IF v_q_count > 0 THEN
        v_base_mark := v_total_marks / v_q_count;
        v_remainder := v_total_marks % v_q_count;

        -- Update questions marks using a CTE to handle remainder distribution
        -- We sort by created_at and ID to ensure deterministic distribution
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

-- Trigger on Exam update (specifically total_marks)
DROP TRIGGER IF EXISTS trigger_redistribute_on_exam_update ON exams;
CREATE TRIGGER trigger_redistribute_on_exam_update
AFTER UPDATE OF total_marks ON exams
FOR EACH ROW
EXECUTE FUNCTION redistribute_exam_marks();

-- Trigger on Question changes
DROP TRIGGER IF EXISTS trigger_redistribute_on_question_change ON questions;
CREATE TRIGGER trigger_redistribute_on_question_change
AFTER INSERT OR DELETE OR UPDATE OF exam_id, deleted_at ON questions
FOR EACH ROW
EXECUTE FUNCTION redistribute_exam_marks();

-- Run redistribution for all existing exams to sync current data
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT id, total_marks FROM exams WHERE deleted_at IS NULL LOOP
        UPDATE questions q
        SET marks = sub.new_mark
        FROM (
            SELECT 
                id,
                (SELECT total_marks FROM exams WHERE id = r.id) / NULLIF(COUNT(*) OVER (), 0) + 
                CASE WHEN ROW_NUMBER() OVER (ORDER BY created_at ASC, id ASC) <= (SELECT total_marks FROM exams WHERE id = r.id) % NULLIF(COUNT(*) OVER (), 0) THEN 1 ELSE 0 END as new_mark
            FROM questions
            WHERE exam_id = r.id AND deleted_at IS NULL
        ) sub
        WHERE q.id = sub.id;
    END LOOP;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_redistribute_on_exam_update ON exams;
DROP TRIGGER IF EXISTS trigger_redistribute_on_question_change ON questions;
DROP FUNCTION IF EXISTS redistribute_exam_marks();
-- +goose StatementEnd
