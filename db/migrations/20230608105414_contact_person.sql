-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE institutions (
    "id" text NOT NULL DEFAULT uuid_generate_v4 (),
    "code" text,
    "name" text,
    "address" text,
    "website" text,
    "postal_code" text,
    "fax" text,
    "telephone" text,
    "business_permit_ojk" text,
    "permit_bursa" text,
    "other_business_permit_ojk" text,
    "type" text,
    "status" text,
    "created_by" text,
    "license" text,
    "operational_status" text,
    "created_at" timestamp,
    "updated_by" text,
    "updated_at" timestamp,
    "deleted_by" text,
    "deleted_at" timestamp,
    "is_deleted" boolean,
    CONSTRAINT company_code_unique_key UNIQUE (code),
    PRIMARY KEY (id)
);
CREATE TABLE institution_division (
    id TEXT NOT NULL DEFAULT uuid_generate_v4(),
    is_default BOOLEAN,
    name TEXT,
    created_by TEXT,
    created_at TIMESTAMPTZ,
    updated_by TEXT,
    updated_at TIMESTAMPTZ,
    deleted_by TEXT,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id)
);
CREATE TABLE institution_members (
    id TEXT NOT NULL DEFAULT uuid_generate_v4(),
    institution_id TEXT NOT NULL,
    division_id TEXT NOT NULL,
    name TEXT NOT NULL,
    position TEXT,
    phone TEXT,
    telephone TEXT,
    email TEXT NOT NULL,
    created_by TEXT,
    created_at TIMESTAMPTZ,
    updated_by TEXT,
    updated_at TIMESTAMPTZ,
    deleted_by TEXT,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    FOREIGN KEY (institution_id) REFERENCES institutions (id),
    FOREIGN KEY (division_id) REFERENCES institution_division (id)
);
INSERT INTO public.institution_division(name, is_default, created_at)
VALUES ('Accounting', true, now()),
    ('IT', true, now()),
    ('Finance', true, now()),
    ('Risk Management', true, now()),
    ('Compliance', true, now()),
    ('Settlement', true, now()),
    ('HR', true, now()),
    ('Marketing', true, now()),
    ('Managerial', true, now());
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DELETE FROM institution_division;
DROP TABLE institution_members;
DROP TABLE institution_division;
DROP TABLE institutions;
-- +goose StatementEnd