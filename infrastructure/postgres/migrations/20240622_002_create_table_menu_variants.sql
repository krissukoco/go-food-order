-- Down
DROP TABLE IF EXISTS "menu_variants";

-- Up
CREATE TABLE IF NOT EXISTS "menu_variants" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "menu_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "min_selected" INT NOT NULL DEFAULT 0,
    "max_selected" INT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);