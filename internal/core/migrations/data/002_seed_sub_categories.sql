BEGIN;

INSERT INTO sub_categories (name, slug, category_id)
VALUES
    ('Телефоны', 'phones',
     (SELECT id FROM categories WHERE slug = 'electronics')),
    ('Ноутбуки', 'laptops',
     (SELECT id FROM categories WHERE slug = 'electronics')),
    ('Мужская', 'men',
     (SELECT id FROM categories WHERE slug = 'clothes')),
    ('Женская', 'women',
     (SELECT id FROM categories WHERE slug = 'clothes')),
    ('Детская', 'child',
    (SELECT id FROM categories WHERE slug = 'clothes'));


COMMIT;

