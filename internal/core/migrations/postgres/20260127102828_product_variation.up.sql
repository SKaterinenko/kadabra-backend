CREATE TABLE IF NOT EXISTS product_variations(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image text not null unique,
    price NUMERIC(10,2) not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER product_variations_set_updated_at
    BEFORE UPDATE ON product_variations
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();