-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE IF EXISTS public.guidance_file_and_regulation;
CREATE TABLE public.guidance_file_and_regulation (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    category text NOT NULL,
    "name" text NOT NULL,
    "description" text,
    link text,
    "file" text,
    file_path text,
    file_size integer,
    "version" "text",
    file_group text,
    "owner" text,
    is_visible boolean,
    created_by text,
    created_at bigint,
    updated_by text,
    updated_at bigint,
    deleted_by text,
    deleted_at bigint,
    PRIMARY KEY (id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS public.guidance_file_and_regulation;
-- +goose StatementEnd