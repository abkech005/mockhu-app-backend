-- Postfeed Likes table
CREATE TABLE postfeed_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    postfeed_id UUID NOT NULL REFERENCES postfeeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(postfeed_id, user_id)
);

-- Postfeed Comments table
CREATE TABLE postfeed_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    postfeed_id UUID NOT NULL REFERENCES postfeeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES postfeed_comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    like_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT comment_content_length CHECK (char_length(content) >= 1 AND char_length(content) <= 2000)
);

-- Postfeed Shares table
CREATE TABLE postfeed_shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    postfeed_id UUID NOT NULL REFERENCES postfeeds(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_postfeed_likes_postfeed ON postfeed_likes(postfeed_id);
CREATE INDEX idx_postfeed_likes_user ON postfeed_likes(user_id);
CREATE INDEX idx_postfeed_comments_postfeed ON postfeed_comments(postfeed_id);
CREATE INDEX idx_postfeed_comments_user ON postfeed_comments(user_id);
CREATE INDEX idx_postfeed_comments_parent ON postfeed_comments(parent_id);
CREATE INDEX idx_postfeed_shares_postfeed ON postfeed_shares(postfeed_id);
CREATE INDEX idx_postfeed_shares_user ON postfeed_shares(user_id);

-- Comments
COMMENT ON TABLE postfeed_likes IS 'User likes on postfeeds';
COMMENT ON TABLE postfeed_comments IS 'User comments on postfeeds with reply support';
COMMENT ON TABLE postfeed_shares IS 'User shares of postfeeds';
