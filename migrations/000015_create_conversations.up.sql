-- Create conversations table for 1-on-1 messaging
CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Participants (always exactly 2 users)
    -- user1_id is always less than user2_id to ensure uniqueness
    user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Last message info (denormalized for performance in inbox listing)
    last_message_id UUID,
    last_message_text TEXT,
    last_message_sender_id UUID,
    last_message_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT unique_conversation UNIQUE (user1_id, user2_id),
    CONSTRAINT different_users CHECK (user1_id != user2_id),
    CONSTRAINT ordered_users CHECK (user1_id < user2_id)
);

-- Indexes for efficient queries
CREATE INDEX idx_conversations_user1 ON conversations(user1_id);
CREATE INDEX idx_conversations_user2 ON conversations(user2_id);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at DESC NULLS LAST);
CREATE INDEX idx_conversations_updated_at ON conversations(updated_at DESC);

-- Add comment
COMMENT ON TABLE conversations IS '1-on-1 conversations between two users';
COMMENT ON COLUMN conversations.user1_id IS 'First user ID (always < user2_id)';
COMMENT ON COLUMN conversations.user2_id IS 'Second user ID (always > user1_id)';
COMMENT ON CONSTRAINT ordered_users ON conversations IS 'Ensures user1_id < user2_id to prevent duplicate conversations';

