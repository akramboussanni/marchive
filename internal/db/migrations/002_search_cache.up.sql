-- Search cache table to store search results temporarily
-- Compatible with both SQLite and PostgreSQL
CREATE TABLE search_cache (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    query TEXT NOT NULL,
    results TEXT NOT NULL, -- JSON encoded search results
    total_results INTEGER NOT NULL,
    created_at BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index for search cache lookup and cleanup
CREATE INDEX idx_search_cache_user_id ON search_cache(user_id);
CREATE INDEX idx_search_cache_expires_at ON search_cache(expires_at);
CREATE INDEX idx_search_cache_created_at ON search_cache(created_at);


