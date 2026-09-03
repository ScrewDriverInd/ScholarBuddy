-- DESTRUCTIVE: clears ScholarBuddy application data and migration history.
-- It does not touch Supabase's auth schema or Google OAuth configuration.
DROP TABLE IF EXISTS opportunities CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS approval_status CASCADE;
DROP TYPE IF EXISTS opportunity_type CASCADE;
DROP TYPE IF EXISTS user_role CASCADE;
DROP TABLE IF EXISTS schema_migrations;
