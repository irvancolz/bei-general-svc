-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE IF EXISTS public.guidance_file_and_regulation
    ADD COLUMN file_size text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE public.guidance_file_and_regulation DROP file_size; 
-- +goose StatementEnd
