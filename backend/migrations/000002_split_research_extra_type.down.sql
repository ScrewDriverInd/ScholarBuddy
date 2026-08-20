-- A rollback intentionally combines research and extras back into the legacy
-- value because the old schema cannot represent them separately.
ALTER TYPE opportunity_type RENAME TO opportunity_type_split;

CREATE TYPE opportunity_type AS ENUM (
    'scholarship',
    'hackathon',
    'internship',
    'research_extra'
);

ALTER TABLE opportunities
    ALTER COLUMN type TYPE opportunity_type
    USING (
        CASE type::text
            WHEN 'research' THEN 'research_extra'
            WHEN 'extras' THEN 'research_extra'
            ELSE type::text
        END
    )::opportunity_type;

DROP TYPE opportunity_type_split;
