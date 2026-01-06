-- Remove index
DROP INDEX IF EXISTS idx_savedbooks_download_count;

-- Remove download_count column from savedbooks table
ALTER TABLE savedbooks DROP COLUMN download_count;
