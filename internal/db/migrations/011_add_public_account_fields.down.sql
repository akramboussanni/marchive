-- Drop anonymous downloads table
DROP INDEX IF EXISTS idx_anonymous_downloads_ip_date;
DROP INDEX IF EXISTS idx_anonymous_downloads_md5;
DROP TABLE IF EXISTS anonymous_downloads;

-- Remove daily download limit from users
ALTER TABLE users DROP COLUMN IF EXISTS daily_download_limit;
