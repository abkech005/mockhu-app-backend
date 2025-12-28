-- Create titles table
CREATE TABLE IF NOT EXISTS titles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    defined_by VARCHAR(10) NOT NULL CHECK (defined_by IN ('admin', 'user')),
    used_by_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_titles_name ON titles(name);
CREATE INDEX IF NOT EXISTS idx_titles_defined_by ON titles(defined_by);
CREATE INDEX IF NOT EXISTS idx_titles_used_by_count ON titles(used_by_count DESC);

-- Add comments for documentation
COMMENT ON TABLE titles IS 'Stores user titles that can be defined by admins or users';
COMMENT ON COLUMN titles.description IS 'Brief description of the title';
COMMENT ON COLUMN titles.defined_by IS 'Who defined the title: admin or user';
COMMENT ON COLUMN titles.used_by_count IS 'Number of users currently using this title';
