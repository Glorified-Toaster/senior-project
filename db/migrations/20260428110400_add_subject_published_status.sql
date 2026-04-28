-- +goose Up
-- +goose StatementBegin
ALTER TYPE subject_status_type ADD VALUE 'PUBLISHED';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- PostgreSQL does not easily support removing enum values. 
-- In a real downgrade, we might need to recreate the type or leave the value in place.
-- For this migration, we will leave it as a no-op manually.
-- +goose StatementEnd
