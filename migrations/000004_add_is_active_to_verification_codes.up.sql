-- Add is_active column to verification_codes table
ALTER TABLE verification_codes ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;

-- Add index for better query performance
CREATE INDEX idx_verification_codes_is_active ON verification_codes(is_active);

-- Set existing codes to inactive if they have been used
UPDATE verification_codes SET is_active = false WHERE used_at IS NOT NULL;

