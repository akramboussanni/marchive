-- Migration 007: Remove password reset token fields
-- Remove unused password reset functionality fields

-- PostgreSQL migration - remove password reset columns
ALTER TABLE users DROP COLUMN IF EXISTS password_reset_token;
ALTER TABLE users DROP COLUMN IF EXISTS password_reset_issuedat;
