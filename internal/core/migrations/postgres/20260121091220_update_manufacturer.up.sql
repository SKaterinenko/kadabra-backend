-- Добавляем колонку category_ids как массив INTEGER
ALTER TABLE manufacturers
    ADD COLUMN category_ids INTEGER[] NOT NULL DEFAULT '{}';

-- Создаем индекс для ускорения поиска по категориям (GIN индекс для массивов)
CREATE INDEX idx_manufacturers_category_ids ON manufacturers USING GIN (category_ids);

-- Добавляем constraint для проверки что массив не пустой
ALTER TABLE manufacturers
    ADD CONSTRAINT manufacturers_category_ids_not_empty
        CHECK (array_length(category_ids, 1) > 0);

-- Добавляем функцию для проверки что все ID категорий существуют
CREATE OR REPLACE FUNCTION check_manufacturer_category_ids()
    RETURNS TRIGGER AS $func$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM unnest(NEW.category_ids) AS cat_id
        WHERE cat_id NOT IN (SELECT id FROM categories)
    ) THEN
        RAISE EXCEPTION 'Invalid category_id in array. All category IDs must exist in categories table.';
    END IF;
    RETURN NEW;
END;
$func$ LANGUAGE plpgsql;

-- Создаем триггер для проверки при INSERT и UPDATE
CREATE TRIGGER trigger_check_manufacturer_category_ids
    BEFORE INSERT OR UPDATE OF category_ids ON manufacturers
    FOR EACH ROW
EXECUTE FUNCTION check_manufacturer_category_ids();