-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create user_roles table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

-- Insert default roles
INSERT INTO roles (id, name, description, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'admin', 'Administrator with full access', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'user', 'Regular user with limited access', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'survey_creator', 'User who can create surveys', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- Insert admin user with password 'admin123' (this is for development only)
INSERT INTO users (id, username, email, password_hash, first_name, last_name, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    'admin',
    'admin@pollpulse.com',
    '$2a$10$XJrL9O7UXs/H5HU0qLXWFesLQQvB2KyPkWPyrXWn.iXJXWE86tdl.', -- 'admin123'
    'Admin',
    'User',
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- Assign admin role to admin user
INSERT INTO user_roles (user_id, role_id, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000000',
    '11111111-1111-1111-1111-111111111111',
    NOW()
) ON CONFLICT (user_id, role_id) DO NOTHING; 