-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE public.announcements
ADD COLUMN company_code text NULL,
    ADD COLUMN form text NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.announcements DROP COLUMN company_code,
    DROP COLUMN form;
-- +goose StatementEnd