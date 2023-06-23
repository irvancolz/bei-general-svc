-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE IF EXISTS public.pkp;
CREATE TABLE public.pkp (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    stakeholders text NOT NULL,
    code text NOT NULL,
    "name" text NOT NULL,
    external_type text,
    additional_information text,
    question_date timestamptz,
    question text NOT NULL,
    answers text NOT NULL,
    answers_by text NOT NULL,
    answers_at timestamptz,
    topic text,
    file_name text,
    file_path text,
    created_by text NOT NULL,
    created_at timestamptz,
    updated_by text,
    updated_at timestamptz,
    deleted_by text,
    deleted_at timestamptz
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.pkp;
-- +goose StatementEnd