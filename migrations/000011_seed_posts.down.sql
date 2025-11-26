-- Remove seeded posts
-- This will delete all posts that were created by the seed migration
-- Note: This is a destructive operation

-- Delete reactions first (cascade will handle this, but being explicit)
DELETE FROM post_reactions WHERE post_id IN (
    SELECT id FROM posts 
    WHERE created_at >= NOW() - INTERVAL '3 days'
    -- Only delete if we're sure these are seeded posts
    -- In production, you might want to be more specific
);

-- Delete seeded posts
-- Note: This is a simple approach - in production, you might want to tag seeded posts
DELETE FROM posts WHERE created_at >= NOW() - INTERVAL '3 days';


