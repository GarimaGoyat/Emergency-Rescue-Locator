-- Run after registering a user to promote them to admin:
-- UPDATE users SET role = 'admin' WHERE email = 'your-admin@email.com';

-- Or create a default admin (password: admin123 — change immediately in production):
-- Password hash is bcrypt of "admin123"
INSERT INTO users (id, email, password, name, phone, role, created_at, updated_at)
VALUES (
    uuid_generate_v4(),
    'admin@rescue.local',
    '$2a$10$3yyGDWyzr3YGK7T3yPNrZ.DQLJKYKtg2B0GNoFXnPL/a2wEJh3tPW',
    'System Admin',
    '+10000000000',
    'admin',
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;
