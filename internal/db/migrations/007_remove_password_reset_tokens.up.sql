-- Migration 007: Remove password reset token fields
-- Remove unused password reset functionality fields
-- SQLite-compatible migration using table recreation

-- Create new users table without password reset columns
CREATE TABLE users_new (
    id BIGINT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    user_role TEXT NOT NULL,
    jwt_session_id BIGINT,
    request_credits INTEGER NOT NULL DEFAULT 0,
    invite_tokens INTEGER NOT NULL DEFAULT 0
);

-- Copy data from old table to new table
INSERT INTO users_new (id, username, password_hash, created_at, user_role, jwt_session_id, request_credits, invite_tokens)
SELECT id, username, password_hash, created_at, user_role, jwt_session_id, request_credits, invite_tokens
FROM users;

-- Drop old table
DROP TABLE users;

-- Rename new table to users
ALTER TABLE users_new RENAME TO users;
