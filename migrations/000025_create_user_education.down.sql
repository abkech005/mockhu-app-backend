-- Remove added columns from user_education table
ALTER TABLE user_education DROP COLUMN IF EXISTS grade;
ALTER TABLE user_education DROP COLUMN IF EXISTS activities;
ALTER TABLE user_education DROP COLUMN IF EXISTS description;

-- Drop added index
DROP INDEX IF EXISTS idx_user_education_current;
