-- Seed Users
-- Passwords are set to 'password123' (plain text for demo, usually should be hashed)
INSERT INTO users (username, email, password_hash)
VALUES 
    ('admin_one', 'admin1@example.com', '$2a$10$8K1p/a0WlE9QvD1p8Y.vEe1zX5zX5zX5zX5zX5zX5zX5zX5zX5zX5'), -- hashed 'password123'
    ('admin_two', 'admin2@example.com', '$2a$10$8K1p/a0WlE9QvD1p8Y.vEe1zX5zX5zX5zX5zX5zX5zX5zX5zX5zX5'),
    ('user_one', 'user1@example.com', '$2a$10$8K1p/a0WlE9QvD1p8Y.vEe1zX5zX5zX5zX5zX5zX5zX5zX5zX5zX5'),
    ('user_two', 'user2@example.com', '$2a$10$8K1p/a0WlE9QvD1p8Y.vEe1zX5zX5zX5zX5zX5zX5zX5zX5zX5zX5')
ON CONFLICT (username) DO NOTHING;

-- Assign Roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin_one' AND r.name = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin_two' AND r.name = 'admin'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'user_one' AND r.name = 'user'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'user_two' AND r.name = 'user'
ON CONFLICT DO NOTHING;
