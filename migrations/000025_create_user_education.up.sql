-- Add new columns to existing user_education table
ALTER TABLE user_education ADD COLUMN IF NOT EXISTS grade VARCHAR(20);
ALTER TABLE user_education ADD COLUMN IF NOT EXISTS activities TEXT;
ALTER TABLE user_education ADD COLUMN IF NOT EXISTS description TEXT;

-- Create additional indexes
CREATE INDEX IF NOT EXISTS idx_user_education_current ON user_education(current);

-- Add comments for documentation
COMMENT ON TABLE user_education IS 'Stores user education history';
COMMENT ON COLUMN user_education.school IS 'Name of school, college, or university';
COMMENT ON COLUMN user_education.degree IS 'Degree type like B.Tech, MBA, PhD';
COMMENT ON COLUMN user_education.field_of_study IS 'Major or specialization';
COMMENT ON COLUMN user_education.current IS 'Whether user is currently enrolled';
COMMENT ON COLUMN user_education.grade IS 'CGPA, percentage, or class obtained';
COMMENT ON COLUMN user_education.activities IS 'Extracurricular activities, achievements';
