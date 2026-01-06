-- Add ghost mode columns to savedbooks
ALTER TABLE savedbooks ADD COLUMN is_ghost BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE savedbooks ADD COLUMN requested_by BIGINT;

-- Add foreign key constraint for requested_by
-- Note: For SQLite, we need to handle this differently if needed
-- For now, we'll add it as a regular column without FK constraint in SQLite
-- PostgreSQL will handle it properly

-- Index for ghost mode filtering
CREATE INDEX idx_savedbooks_is_ghost ON savedbooks(is_ghost);
CREATE INDEX idx_savedbooks_requested_by ON savedbooks(requested_by);
