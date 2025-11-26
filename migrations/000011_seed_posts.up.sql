-- Seed sample posts for testing
-- This migration creates sample posts from existing users
-- Note: This will only work if users already exist in the database

-- Get a few user IDs and create diverse sample posts
DO $$
DECLARE
    user1_id UUID;
    user2_id UUID;
    user3_id UUID;
    post_id UUID;
BEGIN
    -- Get first 3 users (or create test users if none exist)
    SELECT id INTO user1_id FROM users ORDER BY created_at LIMIT 1;
    SELECT id INTO user2_id FROM users ORDER BY created_at OFFSET 1 LIMIT 1;
    SELECT id INTO user3_id FROM users ORDER BY created_at OFFSET 2 LIMIT 1;

    -- Only seed if we have at least one user
    IF user1_id IS NOT NULL THEN
        -- User 1 posts (Technology focused)
        INSERT INTO posts (user_id, content, images, is_anonymous, view_count, created_at) VALUES
        (user1_id, 'Just finished building my first API with Go and Fiber! The performance is incredible. 🚀 #coding #golang', 
         ARRAY['https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=800'], false, 42, NOW() - INTERVAL '2 days'),
        
        (user1_id, 'Working on a new social media platform. The architecture is coming together nicely. Excited to share more soon! 💻', 
         ARRAY[]::TEXT[]::TEXT[], false, 15, NOW() - INTERVAL '1 day'),
        
        (user1_id, 'Anyone else love the feeling when your code finally compiles after hours of debugging? 😅', 
         ARRAY[]::TEXT[]::TEXT[], false, 89, NOW() - INTERVAL '5 hours'),
        
        (user1_id, 'Just discovered PostgreSQL arrays. Mind blown! 🤯 The flexibility is amazing.', 
         ARRAY['https://images.unsplash.com/photo-1544383835-b9af0e3b90f9?w=800'], false, 23, NOW() - INTERVAL '3 hours');

        -- User 2 posts (if exists - Lifestyle/Arts)
        IF user2_id IS NOT NULL THEN
            INSERT INTO posts (user_id, content, images, is_anonymous, view_count, created_at) VALUES
            (user2_id, 'Beautiful sunset today! 🌅 Sometimes you just need to stop and appreciate the little things.', 
             ARRAY['https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800'], false, 67, NOW() - INTERVAL '1 day'),
            
            (user2_id, 'Started learning photography. Here''s my first attempt at street photography. Feedback welcome! 📷', 
             ARRAY['https://images.unsplash.com/photo-1492684223066-81342ee5ff30?w=800', 'https://images.unsplash.com/photo-1502920917128-1aa500764cbd?w=800'], false, 124, NOW() - INTERVAL '12 hours'),
            
            (user2_id, 'Coffee and coding. The perfect morning routine ☕', 
             ARRAY[]::TEXT[], false, 45, NOW() - INTERVAL '6 hours'),
            
            (user2_id, 'Just finished reading an amazing book on design systems. The principles apply everywhere! 📚', 
             ARRAY[]::TEXT[], false, 31, NOW() - INTERVAL '2 hours');
        END IF;

        -- User 3 posts (if exists - Mixed content)
        IF user3_id IS NOT NULL THEN
            INSERT INTO posts (user_id, content, images, is_anonymous, view_count, created_at) VALUES
            (user3_id, 'Working out at the gym today. Consistency is key! 💪', 
             ARRAY[]::TEXT[], false, 56, NOW() - INTERVAL '1 day'),
            
            (user3_id, 'New recipe tried today: Homemade pasta. Turned out amazing! 🍝', 
             ARRAY['https://images.unsplash.com/photo-1551462147-5bc923f49fef?w=800'], false, 78, NOW() - INTERVAL '8 hours'),
            
            (user3_id, 'Travel tip: Always pack a power bank. You''ll thank yourself later! ✈️', 
             ARRAY[]::TEXT[], false, 92, NOW() - INTERVAL '4 hours'),
            
            (user3_id, 'Anonymous post: Sometimes I wonder if anyone actually reads these... 🤔', 
             ARRAY[]::TEXT[], true, 12, NOW() - INTERVAL '1 hour');
        END IF;

        -- If we only have one user, create more diverse posts from that user
        IF user2_id IS NULL AND user3_id IS NULL THEN
            INSERT INTO posts (user_id, content, images, is_anonymous, view_count, created_at) VALUES
            (user1_id, 'Weekend project: Building a REST API with proper error handling. Learning so much! 📖', 
             ARRAY[]::TEXT[], false, 34, NOW() - INTERVAL '10 hours'),
            
            (user1_id, 'Code review tip: Always explain the "why" not just the "what". Makes a huge difference! 💡', 
             ARRAY[]::TEXT[], false, 67, NOW() - INTERVAL '7 hours'),
            
            (user1_id, 'Just deployed my first production API. The feeling is unreal! 🎉', 
             ARRAY['https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800'], false, 156, NOW() - INTERVAL '4 hours'),
            
            (user1_id, 'Database migrations are like time travel for your schema. So powerful! ⏰', 
             ARRAY[]::TEXT[], false, 28, NOW() - INTERVAL '2 hours');
        END IF;
    END IF;
END $$;

-- Add some reactions to make posts more interesting
DO $$
DECLARE
    post_rec RECORD;
    user_rec RECORD;
    reaction_count INT := 0;
BEGIN
    -- Add reactions to random posts
    FOR post_rec IN SELECT id FROM posts ORDER BY created_at DESC LIMIT 10 LOOP
        -- Get a random user (not the post owner)
        FOR user_rec IN 
            SELECT u.id FROM users u 
            WHERE u.id != (SELECT user_id FROM posts WHERE id = post_rec.id)
            ORDER BY RANDOM() 
            LIMIT (FLOOR(RANDOM() * 5)::INT + 1)
        LOOP
            INSERT INTO post_reactions (post_id, user_id, reaction_type, created_at)
            VALUES (post_rec.id, user_rec.id, 'fire', NOW() - (RANDOM() * INTERVAL '2 days'))
            ON CONFLICT (post_id, user_id) DO NOTHING;
            
            reaction_count := reaction_count + 1;
        END LOOP;
    END LOOP;
END $$;

