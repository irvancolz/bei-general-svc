-- +goose Up
-- +goose StatementBegin
-- SELECT 'up SQL query';
CREATE TABLE public.institutes (
	id uuid NOT NULL DEFAULT PRIMARY KEY uuid_generate_v4(),
	external_type text NULL,
	code text NULL,
	institute_name text NULL,
	created_by uuid NULL,
	created_at timestamptz NULL,
	updated_by uuid NULL,
	updated_at timestamptz NULL,
	is_deleted bool NULL,
	deleted_by uuid NULL,
	deleted_at timestamptz NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "institutes";
-- SELECT 'down SQL query';
-- +goose StatementEnd
