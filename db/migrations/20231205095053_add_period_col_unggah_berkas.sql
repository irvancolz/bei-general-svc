-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.uploaded_files
ADD COLUMN periode bigint;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd