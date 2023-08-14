-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.topics ADD COLUMN company_id text NULL, ADD COLUMN user_type text NULL, ADD COLUMN external_type text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.topics DROP COLUMN company_id, DROP COLUMN user_type, DROP COLUMN external_type;
-- +goose StatementEnd
