DROP INDEX opportunities_approval_status_idx;
ALTER TABLE opportunities DROP COLUMN approval_status;

DROP INDEX users_username_idx;
ALTER TABLE users DROP COLUMN password_hash;
ALTER TABLE users DROP COLUMN username;
ALTER TABLE users DROP CONSTRAINT users_roles_not_empty;
ALTER TABLE users RENAME COLUMN roles TO role;
ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users
    ALTER COLUMN role TYPE user_role
    USING role[1];
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'ROLE_USER';
