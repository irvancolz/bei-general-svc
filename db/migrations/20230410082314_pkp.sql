-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.pkp (
id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
stakeholders text NOT NULL,
code text NOT NULL,
name text NOT NULL,
question_date timestamptz WITH TIME ZONE NULL,
question text NOT NULL,
answers text NOT NULL,
answers_by text NOT NULL,
answers_at timestamptz WITH TIME ZONE NOT NULL,
topic text NOT NULL,
file_name text NOT NULL,
file_path text NOT NULL,
create_by text NOT NULL,
created_at timestamptz WITH TIME ZONE NOT NULL,
updated_by text NOT NULL,
updated_at timestamptz WITH TIME ZONE NOT NULL,
deleted_by text,
deleted_at timestamptz WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
DROP TABLE IF EXISTS "pkp";
-- +goose StatementEnd