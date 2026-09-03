-- Drop admins table
DROP TABLE IF EXISTS admins;

-- Restore password_hash column to users table
ALTER TABLE users ADD COLUMN password_hash TEXT;

-- Drop unique constraint on username
DROP INDEX IF EXISTS users_username_unique_idx;

-- Make username nullable again
ALTER TABLE users ALTER COLUMN username DROP NOT NULL;
ALTER TABLE users ALTER COLUMN username DROP DEFAULT;
