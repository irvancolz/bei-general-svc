-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.faqs ADD COLUMN status text DEFAULT 'PUBLISHED';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.faqs DROP COLUMN status;
-- +goose StatementEnd
