DROP TRIGGER IF EXISTS categories_set_updated_at ON categories;
DROP TRIGGER IF EXISTS sub_categories_set_updated_at ON sub_categories;
DROP TRIGGER IF EXISTS products_set_updated_at ON products;
DROP TRIGGER IF EXISTS manufacturers_set_updated_at ON manufacturers;

DROP FUNCTION IF EXISTS set_updated_at();
