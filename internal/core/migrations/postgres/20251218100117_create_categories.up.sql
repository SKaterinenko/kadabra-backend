CREATE TABLE IF NOT EXISTS categories(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS category_translations(
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    language_code VARCHAR(2) NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, language_code),
    UNIQUE(language_code, slug)
);

CREATE TABLE IF NOT EXISTS sub_categories(
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sub_category_translations(
    id SERIAL PRIMARY KEY,
    sub_category_id INTEGER NOT NULL REFERENCES sub_categories(id) ON DELETE CASCADE,
    language_code VARCHAR(2) NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(sub_category_id, language_code)
);

CREATE TABLE IF NOT EXISTS manufacturers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS manufacturer_translations(
    id SERIAL PRIMARY KEY,
    manufacturer_id INTEGER NOT NULL REFERENCES manufacturers(id) ON DELETE CASCADE,
    language_code VARCHAR(2) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(manufacturer_id, language_code)
);

CREATE TABLE IF NOT EXISTS products_type (
    id SERIAL PRIMARY KEY,
    sub_category_id INT NOT NULL REFERENCES sub_categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_type_translations(
    id SERIAL PRIMARY KEY,
    product_type_id INTEGER NOT NULL REFERENCES products_type(id) ON DELETE CASCADE,
    language_code VARCHAR(2) NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_type_id, language_code)
    );

-- Триггер для проверки уникальности
CREATE OR REPLACE FUNCTION check_product_type_name_uniqueness()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM product_type_translations ptt
        JOIN products_type pt ON ptt.product_type_id = pt.id
        JOIN products_type pt2 ON pt2.sub_category_id = pt.sub_category_id
        WHERE pt2.id = NEW.product_type_id
          AND ptt.language_code = NEW.language_code
          AND ptt.name = NEW.name
          AND ptt.product_type_id != NEW.product_type_id
    ) THEN
        RAISE EXCEPTION 'Product type name "%" already exists in this subcategory for language "%"',
            NEW.name, NEW.language_code;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_product_type_name_uniqueness
    BEFORE INSERT OR UPDATE ON product_type_translations
    FOR EACH ROW EXECUTE FUNCTION check_product_type_name_uniqueness();

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_type_id INT NOT NULL REFERENCES products_type(id) ON DELETE CASCADE,
    manufacturer_id INT NOT NULL REFERENCES manufacturers(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_translations(
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    language_code VARCHAR(2) NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    short_description TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, language_code),
    UNIQUE(language_code, slug)
    );

