-- Add new columns to interests table
ALTER TABLE interests ADD COLUMN IF NOT EXISTS defined_by VARCHAR(10) DEFAULT 'admin';
ALTER TABLE interests ADD COLUMN IF NOT EXISTS used_by_count INTEGER DEFAULT 0;
ALTER TABLE interests ADD COLUMN IF NOT EXISTS description VARCHAR(255);

-- Add check constraint for defined_by
ALTER TABLE interests ADD CONSTRAINT interests_defined_by_check CHECK (defined_by IN ('admin', 'user'));

-- Update existing interests to be admin-defined
UPDATE interests SET defined_by = 'admin' WHERE defined_by IS NULL;

-- Calculate and set used_by_count from user_interests table
UPDATE interests i
SET used_by_count = (
    SELECT COUNT(*) 
    FROM user_interests ui 
    WHERE ui.interest_id = i.id
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_interests_defined_by ON interests(defined_by);
CREATE INDEX IF NOT EXISTS idx_interests_used_by_count ON interests(used_by_count DESC);

-- Add comments for documentation
COMMENT ON COLUMN interests.defined_by IS 'Who defined the interest: admin or user';
COMMENT ON COLUMN interests.used_by_count IS 'Number of users who have this interest';
COMMENT ON COLUMN interests.description IS 'Brief description of the interest';
