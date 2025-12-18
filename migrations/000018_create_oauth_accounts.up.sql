CREATE TABLE IF NOT EXISTS oauth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(20) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    
    UNIQUE(provider, provider_user_id),
    UNIQUE(user_id, provider)
);

CREATE INDEX IF NOT EXISTS idx_oauth_provider_user ON oauth_accounts(provider, provider_user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_user_id ON oauth_accounts(user_id);
