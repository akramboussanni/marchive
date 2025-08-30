-- Remove request credits functionality

-- Drop request_credits_log table
DROP TABLE IF EXISTS request_credits_log;

-- Remove request_credits column from users table
ALTER TABLE users DROP COLUMN IF EXISTS request_credits;
