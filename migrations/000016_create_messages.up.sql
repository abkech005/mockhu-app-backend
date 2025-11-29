-- Create messages table for conversation messages
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Conversation reference
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    
    -- Message content
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_type VARCHAR(20) NOT NULL CHECK (message_type IN ('text', 'image', 'file')),
    content TEXT,
    
    -- Media attachments (JSONB for flexible storage)
    -- Stores array of {type, url, filename, size, mime_type, width, height, etc.}
    attachments JSONB,
    
    -- Message status
    status VARCHAR(20) NOT NULL DEFAULT 'sent' CHECK (status IN ('sent', 'delivered', 'read')),
    
    -- Read tracking
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMP,
    
    -- Soft delete (delete for me)
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMP,
    deleted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT attachments_required CHECK (
        (message_type IN ('image', 'file') AND attachments IS NOT NULL) OR 
        (message_type = 'text')
    ),
    CONSTRAINT content_or_attachments CHECK (
        content IS NOT NULL OR attachments IS NOT NULL
    )
);

-- Indexes for efficient queries
-- Composite index for conversation message history (most common query)
CREATE INDEX idx_messages_conversation ON messages(conversation_id, created_at DESC);

-- Index for sender queries
CREATE INDEX idx_messages_sender ON messages(sender_id);

-- Partial index for unread messages (faster unread count)
CREATE INDEX idx_messages_is_read ON messages(is_read) WHERE is_read = FALSE;

-- Index for created_at for chronological queries
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);

-- Index for attachments queries (GIN index for JSONB)
CREATE INDEX idx_messages_attachments ON messages USING GIN(attachments);

-- Add comments
COMMENT ON TABLE messages IS 'Messages within conversations (text, images, files)';
COMMENT ON COLUMN messages.message_type IS 'Type: text, image, or file';
COMMENT ON COLUMN messages.attachments IS 'JSONB array of file metadata for images/files';
COMMENT ON COLUMN messages.is_deleted IS 'Soft delete flag for "delete for me" functionality';

