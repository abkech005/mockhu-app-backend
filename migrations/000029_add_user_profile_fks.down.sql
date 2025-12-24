-- Remove title_id and location_id from users table
DROP INDEX IF EXISTS idx_users_location_id;
DROP INDEX IF EXISTS idx_users_title_id;

ALTER TABLE users DROP COLUMN IF EXISTS location_id;
ALTER TABLE users DROP COLUMN IF EXISTS title_id;
