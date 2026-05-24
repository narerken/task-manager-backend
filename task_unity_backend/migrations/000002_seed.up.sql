INSERT INTO roles (name, department_id)
VALUES
    ('admin',  NULL),
    ('manager', NULL),
    ('user',  NULL)
    ON CONFLICT DO NOTHING;

INSERT INTO departments (name)
VALUES
    ('Engineering'),
    ('Product')
    ON CONFLICT DO NOTHING;