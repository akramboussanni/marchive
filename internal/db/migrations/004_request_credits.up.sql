-- Add request credits table for users to request downloads beyond daily limits
-- Compatible with both SQLite and PostgreSQL

-- Add request_credits column to users table
ALTER TABLE users ADD COLUMN request_credits INTEGER NOT NULL DEFAULT 0;

-- Create request_credits_log table to track credit usage and grants
CREATE TABLE request_credits_log (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    action TEXT NOT NULL, -- 'granted', 'used', 'expired'
    amount INTEGER NOT NULL,
    reason TEXT,
    admin_user_id BIGINT NULL, -- NULL for system actions like usage
    created_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (admin_user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for request_credits_log
CREATE INDEX idx_request_credits_log_user_id ON request_credits_log(user_id);
CREATE INDEX idx_request_credits_log_action ON request_credits_log(action);
CREATE INDEX idx_request_credits_log_created_at ON request_credits_log(created_at);
CREATE INDEX idx_request_credits_log_admin_user_id ON request_credits_log(admin_user_id);
