DROP INDEX IF EXISTS idx_postfeeds_visibility;
ALTER TABLE postfeeds DROP COLUMN IF EXISTS visibility;
