-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE  IF EXISTS public.guidance_file_and_regulation;
CREATE TABLE public.guidance_file_and_regulation
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    category text NOT NULL,
    "name" text NOT NULL,
    "description" text,
    link text,
    "file" text,
    "version" numeric,
    "order" integer,
    created_by text,
    created_at timestamp with time zone,
    updated_by text,
    updated_at timestamp with time zone,
    deleted_by text,
    deleted_at timestamp with time zone,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.guidance_file_and_regulation
    OWNER to devaqi;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE  IF EXISTS public.guidance_file_and_regulation;
-- +goose StatementEnd
