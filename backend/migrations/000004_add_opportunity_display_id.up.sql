-- Keep UUIDs as internal primary keys while exposing a stable, human-friendly
-- numeric identifier to API consumers.
CREATE SEQUENCE opportunities_display_id_seq AS BIGINT;

ALTER TABLE opportunities ADD COLUMN display_id BIGINT;

UPDATE opportunities
SET display_id = nextval('opportunities_display_id_seq')
WHERE display_id IS NULL;

ALTER SEQUENCE opportunities_display_id_seq OWNED BY opportunities.display_id;
ALTER TABLE opportunities
    ALTER COLUMN display_id SET DEFAULT nextval('opportunities_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL;

CREATE UNIQUE INDEX opportunities_display_id_idx ON opportunities (display_id);
