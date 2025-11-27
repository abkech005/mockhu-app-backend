-- Add profile fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS institution_id UUID;

-- Add privacy settings
ALTER TABLE users ADD COLUMN IF NOT EXISTS who_can_message VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN IF NOT EXISTS who_can_see_posts VARCHAR(20) DEFAULT 'everyone';
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_followers_list BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_following_list BOOLEAN DEFAULT true;

-- Add constraints for privacy fields
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'valid_message_privacy') THEN
        ALTER TABLE users ADD CONSTRAINT valid_message_privacy 
            CHECK (who_can_message IN ('everyone', 'followers', 'none'));
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'valid_posts_privacy') THEN
        ALTER TABLE users ADD CONSTRAINT valid_posts_privacy 
            CHECK (who_can_see_posts IN ('everyone', 'followers', 'none'));
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'bio_length_check') THEN
        ALTER TABLE users ADD CONSTRAINT bio_length_check 
            CHECK (LENGTH(bio) <= 500);
    END IF;
END $$;

-- Create unique index on username (case-insensitive)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower 
    ON users(LOWER(username)) WHERE username IS NOT NULL;

-- Add index on institution_id for faster joins
CREATE INDEX IF NOT EXISTS idx_users_institution_id 
    ON users(institution_id) WHERE institution_id IS NOT NULL;