-- Rollback invite system changes

-- Drop invites table
DROP TABLE IF EXISTS invites;

-- Remove invite_tokens column from users table
ALTER TABLE users DROP COLUMN IF EXISTS invite_tokens;
