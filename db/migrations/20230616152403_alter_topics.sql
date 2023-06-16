-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.topics ADD COLUMN handler_name text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.topics DROP COLUMN handler_name;
-- +goose StatementEnd
