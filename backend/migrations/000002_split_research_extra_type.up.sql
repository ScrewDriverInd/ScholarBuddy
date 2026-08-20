-- Replace the legacy combined value with separate public opportunity types.
-- Existing research_extra rows are preserved as research because their original
-- subtype cannot be inferred automatically.
ALTER TYPE opportunity_type RENAME TO opportunity_type_legacy;

CREATE TYPE opportunity_type AS ENUM (
    'scholarship',
    'hackathon',
    'internship',
    'research',
    'extras'
);

ALTER TABLE opportunities
    ALTER COLUMN type TYPE opportunity_type
    USING (
        CASE type::text
            WHEN 'research_extra' THEN 'research'
            ELSE type::text
        END
    )::opportunity_type;

DROP TYPE opportunity_type_legacy;
