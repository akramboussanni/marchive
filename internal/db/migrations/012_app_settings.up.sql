CREATE TABLE app_settings (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- Insert default setting for anonymous access (disabled by default)
INSERT INTO app_settings (key, value, updated_at) VALUES ('anonymous_access_enabled', 'false', EXTRACT(EPOCH FROM NOW())::BIGINT);
