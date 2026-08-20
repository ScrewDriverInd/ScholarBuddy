-- Roles become a set so an administrator can hold both ROLE_USER and ROLE_ADMIN.
ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users
    ALTER COLUMN role TYPE user_role[]
    USING ARRAY[role]::user_role[];
ALTER TABLE users
    ALTER COLUMN role SET DEFAULT ARRAY['ROLE_USER'::user_role];
ALTER TABLE users RENAME COLUMN role TO roles;
ALTER TABLE users ADD CONSTRAINT users_roles_not_empty CHECK (cardinality(roles) > 0);
ALTER TABLE users ADD COLUMN username TEXT;
ALTER TABLE users ADD COLUMN password_hash TEXT;
CREATE UNIQUE INDEX users_username_idx ON users (username) WHERE username IS NOT NULL;

-- New and user-edited opportunities are hidden until an administrator approves them.
ALTER TABLE opportunities
    ADD COLUMN approval_status TEXT NOT NULL DEFAULT 'pending'
    CHECK (approval_status IN ('pending', 'approved'));
CREATE INDEX opportunities_approval_status_idx ON opportunities (approval_status, created_at DESC);
