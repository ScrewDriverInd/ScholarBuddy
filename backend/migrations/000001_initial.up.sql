CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM ('ROLE_USER', 'ROLE_ADMIN');
CREATE TYPE opportunity_type AS ENUM ('scholarship', 'hackathon', 'internship', 'research_extra');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    full_name TEXT NOT NULL DEFAULT '',
    role user_role NOT NULL DEFAULT 'ROLE_USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 250),
    description TEXT NOT NULL CHECK (char_length(description) BETWEEN 1 AND 10000),
    type opportunity_type NOT NULL,
    eligibility TEXT NOT NULL DEFAULT '',
    steps TEXT NOT NULL DEFAULT '',
    benefits TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    referral TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX opportunities_type_created_at_idx ON opportunities (type, created_at DESC);
CREATE INDEX opportunities_created_by_idx ON opportunities (created_by);
