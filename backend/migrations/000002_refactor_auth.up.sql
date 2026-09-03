-- Remove password_hash from users table and auto-generate username from email
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- Add username column as NOT NULL with a generated default
ALTER TABLE users ALTER COLUMN username SET NOT NULL;
ALTER TABLE users ALTER COLUMN username SET DEFAULT '';

-- Generate usernames from email for existing users (extract part before @)
UPDATE users SET username = split_part(email, '@', 1) WHERE username IS NULL OR username = '';

-- Create unique constraint on username
CREATE UNIQUE INDEX IF NOT EXISTS users_username_unique_idx ON users (username);

-- Create admins table for admin authentication with argon2id password hashing
CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create index on username for faster login lookups
CREATE INDEX admins_username_idx ON admins (username);
