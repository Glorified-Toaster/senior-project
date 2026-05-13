-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   -- Only set updated_at to NOW() if it wasn't manually changed in the UPDATE statement
   IF (NEW.updated_at IS NOT DISTINCT FROM OLD.updated_at) THEN
       NEW.updated_at = NOW();
   END IF;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
