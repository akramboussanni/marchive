-- Add redemption codes system for admins to create codes that users can redeem
-- Compatible with both SQLite and PostgreSQL

-- Create redemption_codes table
CREATE TABLE redemption_codes (
    id BIGINT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    invite_tokens INTEGER NOT NULL DEFAULT 0,
    request_credits INTEGER NOT NULL DEFAULT 0,
    max_uses INTEGER NOT NULL DEFAULT 1,
    current_uses INTEGER NOT NULL DEFAULT 0,
    expires_at BIGINT NULL, -- NULL means never expires
    revoked_at BIGINT NULL, -- NULL means not revoked
    created_by BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Create redemption_log table to track code usage
CREATE TABLE redemption_log (
    id BIGINT PRIMARY KEY,
    code_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    redeemed_at BIGINT NOT NULL,
    invite_tokens_granted INTEGER NOT NULL DEFAULT 0,
    request_credits_granted INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (code_id) REFERENCES redemption_codes(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for redemption_codes
CREATE INDEX idx_redemption_codes_code ON redemption_codes(code);
CREATE INDEX idx_redemption_codes_created_by ON redemption_codes(created_by);
CREATE INDEX idx_redemption_codes_expires_at ON redemption_codes(expires_at);
CREATE INDEX idx_redemption_codes_revoked_at ON redemption_codes(revoked_at);
CREATE INDEX idx_redemption_codes_created_at ON redemption_codes(created_at);

-- Indexes for redemption_log
CREATE INDEX idx_redemption_log_code_id ON redemption_log(code_id);
CREATE INDEX idx_redemption_log_user_id ON redemption_log(user_id);
CREATE INDEX idx_redemption_log_redeemed_at ON redemption_log(redeemed_at);

-- Unique constraint to prevent multiple redemptions of the same code by the same user
CREATE UNIQUE INDEX idx_redemption_log_code_user ON redemption_log(code_id, user_id);
