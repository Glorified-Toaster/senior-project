-- +goose Up
-- +goose StatementBegin
ALTER TABLE subject_students ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE subject_instructors ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Create trigger for subject_instructors if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'subject_instructors_updated_at') THEN
        CREATE TRIGGER subject_instructors_updated_at
        BEFORE UPDATE ON subject_instructors
        FOR EACH ROW
        EXECUTE FUNCTION set_updated_at();
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS subject_instructors_updated_at ON subject_instructors;
ALTER TABLE subject_instructors DROP COLUMN IF EXISTS updated_at;
ALTER TABLE subject_students DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
