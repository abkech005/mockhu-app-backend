CREATE TABLE post_shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shared_to_type VARCHAR(20) DEFAULT 'timeline',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT valid_share_type CHECK (shared_to_type IN ('timeline', 'dm', 'external')),
    -- Prevent duplicate shares: a user can only share a post once
    CONSTRAINT unique_user_post_share UNIQUE (post_id, user_id)
);

-- Indexes for performance
CREATE INDEX idx_post_shares_post ON post_shares(post_id);
CREATE INDEX idx_post_shares_user ON post_shares(user_id);
CREATE INDEX idx_post_shares_created ON post_shares(created_at DESC);


