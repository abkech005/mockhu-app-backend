-- Add title_id and location_id to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS title_id UUID REFERENCES titles(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS location_id UUID REFERENCES locations(id);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_title_id ON users(title_id);
CREATE INDEX IF NOT EXISTS idx_users_location_id ON users(location_id);
