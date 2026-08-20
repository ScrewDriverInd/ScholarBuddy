DROP INDEX opportunities_display_id_idx;
ALTER TABLE opportunities DROP COLUMN display_id;
DROP SEQUENCE IF EXISTS opportunities_display_id_seq;
