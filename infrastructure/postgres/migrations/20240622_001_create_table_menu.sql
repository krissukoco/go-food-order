-- Down
DROP TABLE IF EXISTS "menus";

-- Up
CREATE TABLE IF NOT EXISTS "menus" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "restaurant_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "image" JSONB NOT NULL DEFAULT '{}',
    "stock" INT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);