-- Add download_count column to savedbooks table
ALTER TABLE savedbooks ADD COLUMN download_count INTEGER NOT NULL DEFAULT 0;

-- Create index for better query performance
CREATE INDEX idx_savedbooks_download_count ON savedbooks(download_count DESC);
