-- Drop indexes
DROP INDEX IF EXISTS idx_user_interests_interest_id;
DROP INDEX IF EXISTS idx_user_interests_user_id;
DROP INDEX IF EXISTS idx_interests_slug;
DROP INDEX IF EXISTS idx_interests_category;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS user_interests;
DROP TABLE IF EXISTS interests;

