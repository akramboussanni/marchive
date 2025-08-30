-- Add invite system for users to invite others
-- Compatible with both SQLite and PostgreSQL

-- Add invite_tokens column to users table (each user gets 1 invite token)
ALTER TABLE users ADD COLUMN invite_tokens INTEGER NOT NULL DEFAULT 1;

-- Create invites table to track generated invites
CREATE TABLE invites (
    id BIGINT PRIMARY KEY,
    token VARCHAR(64) NOT NULL UNIQUE,
    inviter_id BIGINT NOT NULL,
    invitee_username TEXT NULL, -- NULL until used
    invitee_id BIGINT NULL, -- NULL until used
    used_at BIGINT NULL, -- NULL until used
    revoked_at BIGINT NULL, -- NULL until revoked
    created_at BIGINT NOT NULL,
    FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for invites
CREATE INDEX idx_invites_token ON invites(token);
CREATE INDEX idx_invites_inviter_id ON invites(inviter_id);
CREATE INDEX idx_invites_invitee_id ON invites(invitee_id);
CREATE INDEX idx_invites_created_at ON invites(created_at);
CREATE INDEX idx_invites_used_at ON invites(used_at);
CREATE INDEX idx_invites_revoked_at ON invites(revoked_at);
