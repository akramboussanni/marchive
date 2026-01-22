-- Add daily download limit to users table
ALTER TABLE users ADD COLUMN daily_download_limit INTEGER NOT NULL DEFAULT 10;

-- Create anonymous downloads table for tracking IP-based downloads
CREATE TABLE anonymous_downloads (
    id BIGINT PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    md5 TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

-- Indexes for efficient queries
CREATE INDEX idx_anonymous_downloads_ip_date ON anonymous_downloads(ip_address, created_at);
CREATE INDEX idx_anonymous_downloads_md5 ON anonymous_downloads(md5);
