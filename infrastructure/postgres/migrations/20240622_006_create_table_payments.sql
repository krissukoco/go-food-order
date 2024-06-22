-- Down
DROP TABLE IF EXISTS "payments";

-- Up
CREATE TYPE payment_status AS ENUM (
    'PENDING', 'SUCCESS', 'EXPIRED'
);

CREATE TABLE IF NOT EXISTS "payments" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "order_id" UUID NOT NULL,
    "customer_id" UUID NOT NULL,
    "vendor" TEXT NOT NULL,
    "method" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "status" payment_status NOT NULL,
    "expired_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);