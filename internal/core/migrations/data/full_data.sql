-- Electronics / Электроника
WITH new_category AS (
    INSERT INTO categories DEFAULT VALUES RETURNING id
)
INSERT INTO category_translations (category_id, language_code, name, slug)
SELECT id, 'ru', 'Электроника', 'elektronika' FROM new_category
UNION ALL
SELECT id, 'en', 'Electronics', 'electronics' FROM new_category;

-- Clothing and accessories / Одежда и аксессуары
WITH new_category AS (
    INSERT INTO categories DEFAULT VALUES RETURNING id
)
INSERT INTO category_translations (category_id, language_code, name, slug)
SELECT id, 'ru', 'Одежда и аксессуары', 'odezhda-i-aksessuary' FROM new_category
UNION ALL
SELECT id, 'en', 'Clothing and accessories', 'clothing-and-accessories' FROM new_category;

-- Health / Здоровье
WITH new_category AS (
    INSERT INTO categories DEFAULT VALUES RETURNING id
)
INSERT INTO category_translations (category_id, language_code, name, slug)
SELECT id, 'ru', 'Здоровье', 'zdorove' FROM new_category
UNION ALL
SELECT id, 'en', 'Health', 'health' FROM new_category;

-- Cosmetics / Косметика
WITH new_category AS (
    INSERT INTO categories DEFAULT VALUES RETURNING id
)
INSERT INTO category_translations (category_id, language_code, name, slug)
SELECT id, 'ru', 'Косметика', 'kosmetika' FROM new_category
UNION ALL
SELECT id, 'en', 'Cosmetics', 'cosmetics' FROM new_category;

-- House / Дом
WITH new_category AS (
    INSERT INTO categories DEFAULT VALUES RETURNING id
)
INSERT INTO category_translations (category_id, language_code, name, slug)
SELECT id, 'ru', 'Дом', 'dom' FROM new_category
UNION ALL
SELECT id, 'en', 'House', 'house' FROM new_category;