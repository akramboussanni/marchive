-- Initial database schema for marchive
-- All tables and indexes in one migration file

-- Users table with authentication and profile data
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    user_role TEXT NOT NULL,
    -- Password reset fields
    password_reset_token VARCHAR(64),
    password_reset_issuedat BIGINT,
    -- JWT session management
    jwt_session_id BIGINT
);

-- JWT blacklist for token revocation
CREATE TABLE jwt_blacklist (
    jti VARCHAR(255) PRIMARY KEY,
    user_id BIGINT,
    expires_at BIGINT NOT NULL
);

-- Failed login attempts tracking
CREATE TABLE failed_logins (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NULL,
    ip_address VARCHAR(45) NOT NULL,
    attempted_at BIGINT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true
);

-- Account lockouts
CREATE TABLE lockouts (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NULL,
    ip_address VARCHAR(45) NULL,
    locked_until BIGINT NOT NULL,
    reason VARCHAR(255) NULL,
    active BOOLEAN NOT NULL DEFAULT true
);

-- Download rate limiting
CREATE TABLE downloadrequests (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    md5 TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Saved books catalog
CREATE TABLE savedbooks (
    id BIGINT PRIMARY KEY,
    hash TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    authors TEXT,
    publisher TEXT,
    language TEXT,
    format TEXT,
    size TEXT,
    cover_url TEXT,
    cover_data TEXT,
    file_path TEXT,
    status TEXT NOT NULL DEFAULT 'processing',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Download job queue and status
CREATE TABLE downloadjobs (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    book_hash TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    progress INTEGER NOT NULL DEFAULT 0,
    error_msg TEXT,
    file_path TEXT,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (book_hash) REFERENCES savedbooks(hash) ON DELETE CASCADE
);

-- Indexes for failed_logins
CREATE INDEX idx_failed_logins_user ON failed_logins(user_id);
CREATE INDEX idx_failed_logins_ip ON failed_logins(ip_address);
CREATE INDEX idx_failed_logins_attempted_at ON failed_logins(attempted_at);

-- Indexes for lockouts
CREATE INDEX idx_lockouts_user ON lockouts(user_id);
CREATE INDEX idx_lockouts_ip ON lockouts(ip_address);
CREATE INDEX idx_lockouts_locked_until ON lockouts(locked_until);

-- Indexes for downloadrequests
CREATE INDEX idx_downloadrequests_user_date ON downloadrequests(user_id, DATE(created_at, 'unixepoch'));
CREATE INDEX idx_downloadrequests_md5 ON downloadrequests(md5);
CREATE INDEX idx_downloadrequests_user ON downloadrequests(user_id);

-- Indexes for savedbooks
CREATE INDEX idx_savedbooks_hash ON savedbooks(hash);
CREATE INDEX idx_savedbooks_status ON savedbooks(status);
CREATE INDEX idx_savedbooks_title ON savedbooks(title);
CREATE INDEX idx_savedbooks_authors ON savedbooks(authors);

-- Indexes for downloadjobs
CREATE INDEX idx_downloadjobs_user_id ON downloadjobs(user_id);
CREATE INDEX idx_downloadjobs_book_hash ON downloadjobs(book_hash);
CREATE INDEX idx_downloadjobs_status ON downloadjobs(status);
CREATE INDEX idx_downloadjobs_created_at ON downloadjobs(created_at);
