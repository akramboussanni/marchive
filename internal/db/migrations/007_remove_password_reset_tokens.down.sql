-- Migration 007 Down: Restore password reset token fields
-- Add back password reset functionality fields

-- SQLite doesn't support ADD COLUMN easily, so we need to recreate the table
-- First, create a new table with the password reset columns
CREATE TABLE users_old (
    id INTEGER PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    user_role VARCHAR(20) NOT NULL DEFAULT 'user',
    jwt_session_id BIGINT NOT NULL DEFAULT 0,
    request_credits INTEGER NOT NULL DEFAULT 0,
    invite_tokens INTEGER NOT NULL DEFAULT 1,
    password_reset_token VARCHAR(64),
    password_reset_issuedat BIGINT
);

-- Copy data from current table to old table
INSERT INTO users_old (id, username, password_hash, created_at, user_role, jwt_session_id, request_credits, invite_tokens)
SELECT id, username, password_hash, created_at, user_role, jwt_session_id, request_credits, invite_tokens
FROM users;

-- Drop the current table
DROP TABLE users;

-- Rename the old table to the original name
ALTER TABLE users_old RENAME TO users;
