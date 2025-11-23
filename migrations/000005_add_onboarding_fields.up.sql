-- Add onboarding tracking columns to users table
ALTER TABLE users 
  ADD COLUMN IF NOT EXISTS onboarding_completed BOOLEAN DEFAULT false,
  ADD COLUMN IF NOT EXISTS onboarded_at TIMESTAMP WITH TIME ZONE;

-- Add index for faster queries on onboarding status
CREATE INDEX IF NOT EXISTS idx_users_onboarding_completed 
  ON users(onboarding_completed);

-- Set existing users to not onboarded (defensive)
UPDATE users 
  SET onboarding_completed = false 
  WHERE onboarding_completed IS NULL;