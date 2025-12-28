-- Create locations table
CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    used_by_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_city_country UNIQUE (city, country)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_locations_city ON locations(city);
CREATE INDEX IF NOT EXISTS idx_locations_country ON locations(country);
CREATE INDEX IF NOT EXISTS idx_locations_used_by_count ON locations(used_by_count DESC);

-- Add comments
COMMENT ON TABLE locations IS 'Stores city/country locations for users';
COMMENT ON COLUMN locations.used_by_count IS 'Number of users in this location';
