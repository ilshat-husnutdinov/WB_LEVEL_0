CREATE TABLE IF NOT EXISTS orders (
    id        SERIAL PRIMARY KEY,
    order_uid  VARCHAR(100) NOT NULL UNIQUE,
    data JSONB NOT NULL);

CREATE INDEX IF NOT EXISTS idx_oid ON orders(order_uid);