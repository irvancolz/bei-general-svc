-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.faqs ADD COLUMN order_num int DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.faqs DROP COLUMN order_num;
-- +goose StatementEnd
