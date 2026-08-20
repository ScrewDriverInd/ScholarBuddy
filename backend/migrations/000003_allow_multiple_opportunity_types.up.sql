-- Preserve each existing scalar type as a one-item array, then allow an
-- opportunity to belong to any non-empty combination of opportunity types.
DROP INDEX opportunities_type_created_at_idx;

ALTER TABLE opportunities
    ALTER COLUMN type TYPE opportunity_type[]
    USING ARRAY[type]::opportunity_type[];

ALTER TABLE opportunities RENAME COLUMN type TO types;

ALTER TABLE opportunities
    ADD CONSTRAINT opportunities_types_not_empty CHECK (cardinality(types) > 0);

CREATE INDEX opportunities_types_idx ON opportunities USING GIN (types);
CREATE INDEX opportunities_created_at_idx ON opportunities (created_at DESC);
