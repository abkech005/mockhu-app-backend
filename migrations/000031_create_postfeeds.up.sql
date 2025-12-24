-- Postfeeds table - Social media feed with different post types
CREATE TABLE postfeeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Type discriminator: doubt, quiz, progress, resource
    type VARCHAR(20) NOT NULL CHECK (type IN ('doubt', 'quiz', 'progress', 'resource')),
    
    -- Common fields
    title VARCHAR(255) NOT NULL,
    content TEXT,
    tags TEXT[],
    is_anonymous BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    
    -- Type-specific data stored as JSONB
    metadata JSONB DEFAULT '{}',
    
    -- Engagement metrics (denormalized for performance)
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    share_count INTEGER DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX idx_postfeeds_user_id ON postfeeds(user_id);
CREATE INDEX idx_postfeeds_type ON postfeeds(type);
CREATE INDEX idx_postfeeds_created_at ON postfeeds(created_at DESC) WHERE is_active = true;
CREATE INDEX idx_postfeeds_tags ON postfeeds USING GIN(tags);
CREATE INDEX idx_postfeeds_metadata ON postfeeds USING GIN(metadata);

-- Comments
COMMENT ON TABLE postfeeds IS 'Social media feed posts: doubt, quiz, progress, resource';
COMMENT ON COLUMN postfeeds.type IS 'Post type: doubt, quiz, progress, resource';
COMMENT ON COLUMN postfeeds.metadata IS 'Type-specific data in JSONB format';
