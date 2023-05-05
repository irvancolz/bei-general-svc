-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.announcements (
	id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	"type" text NULL,
	information_type text NULL,
	effective_date timestamptz NULL,
	regarding text NULL,
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
-- SELECT 'down SQL query';
DROP TABLE "announcements";
-- +goose StatementEnd
