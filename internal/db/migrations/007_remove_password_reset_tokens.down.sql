-- Migration 007 Down: Restore password reset token fields
-- Add back password reset functionality fields

-- PostgreSQL migration - add password reset columns back
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_reset_token VARCHAR(64);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_reset_issuedat BIGINT;
