DROP TABLE IF EXISTS product_translations;
DROP TABLE IF EXISTS products;
DROP TRIGGER IF EXISTS trigger_check_product_type_name_uniqueness
    ON product_type_translations;
DROP FUNCTION IF EXISTS check_product_type_name_uniqueness();
DROP TABLE IF EXISTS product_type_translations;
DROP TABLE IF EXISTS products_type;
DROP TABLE IF EXISTS manufacturer_translations;
DROP TABLE IF EXISTS manufacturers;
DROP TABLE IF EXISTS sub_category_translations;
DROP TABLE IF EXISTS sub_categories;
DROP TABLE IF EXISTS category_translations;
DROP TABLE IF EXISTS categories;