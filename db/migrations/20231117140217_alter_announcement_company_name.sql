-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.announcements
ADD COLUMN company_name text NULL DEFAULT '';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.announcements DROP COLUMN company_name;
-- +goose StatementEnd