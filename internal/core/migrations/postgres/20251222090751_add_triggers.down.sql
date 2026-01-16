DROP TRIGGER IF EXISTS categories_set_updated_at ON categories;
DROP TRIGGER IF EXISTS sub_categories_set_updated_at ON sub_categories;
DROP TRIGGER IF EXISTS products_type_set_updated_at ON products_type;
DROP TRIGGER IF EXISTS products_set_updated_at ON products;
DROP TRIGGER IF EXISTS manufacturers_set_updated_at ON manufacturers;


DROP TRIGGER IF EXISTS category_translations_set_updated_at ON category_translations;
DROP TRIGGER IF EXISTS manufacturer_translations_set_updated_at ON manufacturer_translations;
DROP TRIGGER IF EXISTS product_translations_set_updated_at ON product_translations;
DROP TRIGGER IF EXISTS product_type_translations_set_updated_at ON product_type_translations;
DROP TRIGGER IF EXISTS sub_category_translations_set_updated_at ON sub_category_translations;

DROP FUNCTION IF EXISTS set_updated_at();
