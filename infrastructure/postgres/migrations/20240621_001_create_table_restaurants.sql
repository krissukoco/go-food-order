-- Down
DROP TABLE IF EXISTS "restaurants";

-- Up
CREATE TABLE IF NOT EXISTS "restaurants" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "name" TEXT NOT NULL,
    "image" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOw()
);