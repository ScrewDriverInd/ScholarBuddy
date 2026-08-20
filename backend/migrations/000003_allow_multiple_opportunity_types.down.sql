-- Rollback retains the first type only because the legacy schema permits one.
DROP INDEX opportunities_types_idx;
DROP INDEX opportunities_created_at_idx;

ALTER TABLE opportunities DROP CONSTRAINT opportunities_types_not_empty;
ALTER TABLE opportunities RENAME COLUMN types TO type;

ALTER TABLE opportunities
    ALTER COLUMN type TYPE opportunity_type
    USING type[1];

CREATE INDEX opportunities_type_created_at_idx ON opportunities (type, created_at DESC);
