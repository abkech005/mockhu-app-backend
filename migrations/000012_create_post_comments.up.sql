CREATE TABLE post_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES post_comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_anonymous BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT content_length CHECK (
        char_length(content) >= 1 AND 
        char_length(content) <= 2000
    )
);

-- Indexes for performance
CREATE INDEX idx_post_comments_post ON post_comments(post_id, created_at);
CREATE INDEX idx_post_comments_user ON post_comments(user_id);
CREATE INDEX idx_post_comments_parent ON post_comments(parent_comment_id);
CREATE INDEX idx_post_comments_active ON post_comments(is_active, created_at DESC);

