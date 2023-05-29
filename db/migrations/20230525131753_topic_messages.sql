-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.topic_messages (
        id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
        topic_id uuid NOT NULL,
        message text NOT NULL,
        company_name text NOT NULL,
        company_id uuid NOT NULL,
        user_full_name text NOT NULL,
        created_by uuid NULL,
        created_at timestamptz NULL,
        updated_by uuid NULL,
        updated_at timestamptz NULL,
        is_deleted bool DEFAULT FALSE,
        deleted_by uuid NULL,
        deleted_at timestamptz NULL
);

ALTER TABLE topic_messages
    ADD CONSTRAINT fk_topic_messages FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE topic_messages
    DROP CONSTRAINT fk_topic_messages;

DROP TABLE "public"."topic_messages";
-- +goose StatementEnd
