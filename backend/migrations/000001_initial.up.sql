CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM ('ROLE_USER', 'ROLE_ADMIN');
CREATE TYPE opportunity_type AS ENUM ('scholarship', 'hackathon', 'internship', 'research', 'extras');
CREATE TYPE approval_status AS ENUM ('pending', 'approved');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL GENERATED ALWAYS AS (split_part(email, '@', 1)) STORED,
    roles user_role[] NOT NULL DEFAULT ARRAY['ROLE_USER'::user_role],
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_roles_not_empty CHECK (cardinality(roles) > 0)
);
CREATE UNIQUE INDEX users_username_idx ON users (username);

CREATE TABLE opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 250),
    description TEXT NOT NULL CHECK (char_length(description) BETWEEN 1 AND 10000),
    types opportunity_type[] NOT NULL CHECK (cardinality(types) > 0),
    eligibility TEXT NOT NULL DEFAULT '',
    steps TEXT NOT NULL DEFAULT '',
    benefits TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    referral TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id),
    approval_status approval_status NOT NULL DEFAULT 'pending',
    click_count BIGINT NOT NULL DEFAULT 0 CHECK (click_count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX opportunities_types_idx ON opportunities USING GIN (types);
CREATE INDEX opportunities_public_rank_idx ON opportunities (approval_status, click_count DESC, created_at DESC);
