-- Down
DROP TABLE IF EXISTS "users_backoffice";

-- Up
CREATE TABLE IF NOT EXISTS "users_backoffice" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "restaurant_id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "password" TEXT NOT NULL,
    "first_password" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE("email")
);