-- Down
DROP TABLE IF EXISTS "customers";

-- Up
CREATE TABLE IF NOT EXISTS "customers" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "phone" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "image" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE("phone")
);