CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER categories_set_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER sub_categories_set_updated_at
    BEFORE UPDATE ON sub_categories
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER products_type_set_updated_at
    BEFORE UPDATE ON products_type
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER products_set_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER manufacturers_set_updated_at
    BEFORE UPDATE ON manufacturers
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER category_translations_set_updated_at
    BEFORE UPDATE ON category_translations
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER manufacturer_translations_set_updated_at
    BEFORE UPDATE ON manufacturer_translations
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER product_translations_set_updated_at
    BEFORE UPDATE ON product_translations
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER product_type_translations_set_updated_at
    BEFORE UPDATE ON product_type_translations
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER sub_category_translations_set_updated_at
    BEFORE UPDATE ON sub_category_translations
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
