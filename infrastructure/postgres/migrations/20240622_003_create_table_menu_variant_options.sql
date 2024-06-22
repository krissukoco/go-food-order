-- Down
DROP TABLE IF EXISTS "menu_variant_options";

-- Up
CREATE TABLE IF NOT EXISTS "menu_variant_options" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "menu_variant_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "additional_cost" INT NOT NULL DEFAULT 0,
    "stock" INT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);