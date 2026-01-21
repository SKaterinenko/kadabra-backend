DROP TRIGGER IF EXISTS trigger_check_manufacturer_category_ids ON manufacturers;
DROP FUNCTION IF EXISTS check_manufacturer_category_ids();
ALTER TABLE manufacturers DROP CONSTRAINT IF EXISTS manufacturers_category_ids_not_empty;
DROP INDEX IF EXISTS idx_manufacturers_category_ids;
ALTER TABLE manufacturers DROP COLUMN IF EXISTS category_ids;