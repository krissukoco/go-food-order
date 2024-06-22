-- Down
DROP TABLE IF EXISTS "order_fulfillment_histories";

-- Up
CREATE TABLE IF NOT EXISTS "order_fulfillment_histories" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "order_id" UUID NOT NULL,
    "status" order_fulfillment_status NOT NULL,
    "action_by" UUID,
    "created_by" UUID NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);