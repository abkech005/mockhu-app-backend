-- Remove onboarding fields from users table
DROP INDEX IF EXISTS idx_users_onboarding_completed;

    
ALTER TABLE users 
  DROP COLUMN IF EXISTS onboarded_at,
  DROP COLUMN IF EXISTS onboarding_completed;