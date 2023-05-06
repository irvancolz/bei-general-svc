-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.pkp (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    stakeholders TEXT NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    question_date TIMESTAMP WITH TIME ZONE NOT NULL,
    question TEXT NOT NULL,
    answers TEXT NOT NULL,
    answers_by TEXT NOT NULL,
    answers_at TIMESTAMP WITH TIME ZONE NOT NULL,
    topic TEXT NOT NULL,
    file_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    create_by TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_by TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_by TEXT,
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
DROP TABLE "pkp";
-- +goose StatementEnd
