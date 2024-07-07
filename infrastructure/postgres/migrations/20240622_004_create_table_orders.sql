-- Down
DROP TABLE IF EXISTS "orders";

-- Up
CREATE TYPE order_payment_method AS ENUM (
    'QRIS', 'OVO', 'SHOPEEPAY', 'VA_MANDIRI', 'VA_BCA', 'VA_BNI'
);

CREATE TYPE order_payment_status AS ENUM (
    'PAID', 'UNPAID'
);

CREATE TYPE order_fulfillment_status AS ENUM (
    'ORDERED', 'CANCELLED', 'CONFIRMED', 'PREPARING', 'PREPARED', 'SERVED'
);

CREATE TYPE order_platform AS ENUM (
    'CUSTOMER_APP', 'POS'
);

CREATE TABLE IF NOT EXISTS "orders" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    "customer_id" UUID NOT NULL,
    "restaurant_id" UUID NOT NULL,
    "table" VARCHAR(16) NOT NULL DEFAULT '',
    "items_count" INT NOT NULL,
    "items_amount" INT NOT NULL,
    "restaurant_fee" INT NOT NULL,
    "service_fee" INT NOT NULL,
    "payment_fee" INT NOT NULL,
    "total_payment" INT NOT NULL,
    "payment_method" order_payment_method NOT NULL,
    "payment_status" order_payment_status NOT NULL,
    "fulfillment_status" order_fulfillment_status NOT NULL,
    "platform" order_platform NOT NULL,
    "notes" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);