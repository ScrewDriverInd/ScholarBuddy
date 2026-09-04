CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM ('ROLE_USER', 'ROLE_ADMIN');
CREATE TYPE opportunity_type AS ENUM ('SCHOLARSHIP', 'HACKATHON', 'INTERNSHIP', 'RESEARCH', 'EXTRAS');
CREATE TYPE approval_status AS ENUM ('PENDING', 'APPROVED');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL DEFAULT '',
    username TEXT NOT NULL GENERATED ALWAYS AS (split_part(email, '@', 1)) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_username_idx ON users (username);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role user_role NOT NULL,
    PRIMARY KEY (user_id, role)
);

CREATE TABLE opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 250),
    description TEXT NOT NULL CHECK (char_length(description) BETWEEN 1 AND 10000),
    eligibility TEXT NOT NULL DEFAULT '',
    steps TEXT NOT NULL DEFAULT '',
    benefits TEXT NOT NULL DEFAULT '',
    link TEXT NOT NULL DEFAULT '',
    referral TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id),
    approval_status approval_status NOT NULL DEFAULT 'PENDING',
    click_count BIGINT NOT NULL DEFAULT 0 CHECK (click_count >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE opportunity_types (
    opportunity_id UUID NOT NULL REFERENCES opportunities(id) ON DELETE CASCADE,
    type opportunity_type NOT NULL,
    PRIMARY KEY (opportunity_id, type)
);

CREATE INDEX opportunities_types_idx ON opportunity_types (type);
CREATE INDEX opportunities_public_rank_idx ON opportunities (approval_status, click_count DESC, created_at DESC);
