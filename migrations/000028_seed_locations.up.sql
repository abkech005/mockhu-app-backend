-- Seed popular locations
INSERT INTO locations (city, country, used_by_count) VALUES
    -- India - Major Cities
    ('Mumbai', 'India', 0),
    ('Delhi', 'India', 0),
    ('Bangalore', 'India', 0),
    ('Hyderabad', 'India', 0),
    ('Chennai', 'India', 0),
    ('Kolkata', 'India', 0),
    ('Pune', 'India', 0),
    ('Ahmedabad', 'India', 0),
    ('Jaipur', 'India', 0),
    ('Lucknow', 'India', 0),
    ('Chandigarh', 'India', 0),
    ('Indore', 'India', 0),
    ('Bhopal', 'India', 0),
    ('Nagpur', 'India', 0),
    ('Kochi', 'India', 0),
    ('Coimbatore', 'India', 0),
    ('Noida', 'India', 0),
    ('Gurugram', 'India', 0),
    ('Surat', 'India', 0),
    ('Vadodara', 'India', 0),
    
    -- International - Major Cities
    ('New York', 'United States', 0),
    ('San Francisco', 'United States', 0),
    ('Los Angeles', 'United States', 0),
    ('Chicago', 'United States', 0),
    ('Seattle', 'United States', 0),
    ('London', 'United Kingdom', 0),
    ('Dubai', 'United Arab Emirates', 0),
    ('Singapore', 'Singapore', 0),
    ('Toronto', 'Canada', 0),
    ('Sydney', 'Australia', 0),
    ('Berlin', 'Germany', 0),
    ('Amsterdam', 'Netherlands', 0),
    ('Tokyo', 'Japan', 0)
ON CONFLICT (city, country) DO NOTHING;
