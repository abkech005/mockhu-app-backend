CREATE TABLE post_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction_type VARCHAR(20) DEFAULT 'fire',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- One reaction per user per post
    UNIQUE(post_id, user_id)
);

-- Indexes
CREATE INDEX idx_post_reactions_post ON post_reactions(post_id);
CREATE INDEX idx_post_reactions_user ON post_reactions(user_id);
CREATE UNIQUE INDEX idx_post_reactions_unique 
    ON post_reactions(post_id, user_id);