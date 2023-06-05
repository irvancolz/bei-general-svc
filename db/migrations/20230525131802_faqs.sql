-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.faqs (
        id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
        question text NOT NULL,
        answer text NOT NULL,
        created_by uuid NULL,
        created_at timestamptz NULL,
        updated_by uuid NULL,
        updated_at timestamptz NULL,
        is_deleted bool DEFAULT FALSE,
        deleted_by uuid NULL,
        deleted_at timestamptz NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "public"."faqs";
-- +goose StatementEnd
