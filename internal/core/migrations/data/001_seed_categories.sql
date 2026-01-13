BEGIN;

INSERT INTO categories (name, slug)
VALUES
    ('Электроника', 'electronics'),
    ('Одежда и аксессуары', 'clothes'),
    ('Здоровье', 'health'),
    ('Косметика', 'beauty'),
    ('Дом', 'home');

COMMIT;
