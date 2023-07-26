-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.log_systems RENAME COLUMN menu TO modul;

ALTER TABLE public.log_systems ADD COLUMN sub_modul text NULL, ADD COLUMN detail text NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.log_systems RENAME COLUMN modul TO menu;

ALTER TABLE public.log_systems DROP COLUMN sub_modul, DROP COLUMN detail;
-- +goose StatementEnd
