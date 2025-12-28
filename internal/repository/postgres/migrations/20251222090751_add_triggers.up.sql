-- 1. Функция (одна на всю БД)
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Триггеры на таблицы
CREATE TRIGGER categories_set_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER sub_categories_set_updated_at
    BEFORE UPDATE ON sub_categories
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER products_set_updated_at
    BEFORE UPDATE ON products    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER manufacturers_set_updated_at
    BEFORE UPDATE ON manufacturers
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
