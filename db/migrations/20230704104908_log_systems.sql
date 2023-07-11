-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.log_systems (
        id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
        menu text NOT NULL,
        action text NOT NULL,
        user_name text NOT NULL,
        ip text NOT NULL,
        browser text NOT NULL,
        created_by uuid NULL,
        created_at timestamptz NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "public"."log_systems";
-- +goose StatementEnd
