-- Drop indexes
DROP INDEX IF EXISTS idx_users_username_lower;
DROP INDEX IF EXISTS idx_users_institution_id;

-- Drop constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS valid_message_privacy;
ALTER TABLE users DROP CONSTRAINT IF EXISTS valid_posts_privacy;
ALTER TABLE users DROP CONSTRAINT IF EXISTS bio_length_check;

-- Drop columns
ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS institution_id;
ALTER TABLE users DROP COLUMN IF EXISTS who_can_message;
ALTER TABLE users DROP COLUMN IF EXISTS who_can_see_posts;
ALTER TABLE users DROP COLUMN IF EXISTS show_followers_list;
ALTER TABLE users DROP COLUMN IF EXISTS show_following_list;