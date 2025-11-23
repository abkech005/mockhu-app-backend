CREATE TABLE IF NOT EXISTS verification_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  code TEXT NOT NULL,
  type TEXT NOT NULL,
  contact TEXT NOT NULL,
  used_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- -- Indexing for performance
CREATE INDEX idx_verification_codes_code_type ON verification_codes(code, type);
CREATE INDEX idx_verification_codes_contact_type ON verification_codes(contact, type);
CREATE INDEX idx_verification_codes_expires_at ON verification_codes(expires_at);
CREATE INDEX idx_verification_codes_user_id ON verification_codes(user_id);