-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.uploaded_files (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    report_type text,
    report_code text,
    report_name text,
    is_uploaded bool,
    file_name text,
    file_path text,
    file_size integer,
    created_by text,
    created_at integer,
    updated_by text,
    updated_at bigint,
    deleted_by text,
    deleted_at bigint
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE uploaded_files;
-- +goose StatementEnd