-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.topics ADD COLUMN company_code text DEFAULT 'BEI', ADD COLUMN company_name text DEFAULT 'Bursa Efek Indonesia';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.topics DROP COLUMN company_code, DROP COLUMN company_name;
-- +goose StatementEnd
