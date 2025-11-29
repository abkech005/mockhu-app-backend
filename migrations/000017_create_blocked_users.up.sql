-- Create blocked_users table for user blocking functionality
CREATE TABLE IF NOT EXISTS blocked_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Who blocked whom
    blocker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Metadata
    reason TEXT,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE(blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);

-- Indexes for efficient blocking checks
CREATE INDEX idx_blocked_users_blocker ON blocked_users(blocker_id);
CREATE INDEX idx_blocked_users_blocked ON blocked_users(blocked_id);

-- Composite index for quick blocking checks
CREATE INDEX idx_blocked_users_both ON blocked_users(blocker_id, blocked_id);

-- Add comments
COMMENT ON TABLE blocked_users IS 'User blocking system - prevents messaging';
COMMENT ON COLUMN blocked_users.blocker_id IS 'User who initiated the block';
COMMENT ON COLUMN blocked_users.blocked_id IS 'User who is blocked';
COMMENT ON COLUMN blocked_users.reason IS 'Optional reason for blocking (for future admin review)';

