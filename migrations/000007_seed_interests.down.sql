-- Remove all seeded interests
-- This will also cascade delete all user_interests due to ON DELETE CASCADE
TRUNCATE interests CASCADE;

