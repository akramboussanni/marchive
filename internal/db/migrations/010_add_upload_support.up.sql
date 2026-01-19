-- Add upload support to savedbooks table
-- This allows user-uploaded books to coexist with downloaded books

-- Add columns for upload tracking
ALTER TABLE savedbooks ADD COLUMN is_uploaded BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE savedbooks ADD COLUMN uploaded_by BIGINT;
ALTER TABLE savedbooks ADD COLUMN original_filename TEXT;

-- Note: Foreign key constraint for uploaded_by -> users(id) is enforced at application level
-- SQLite doesn't support ADD CONSTRAINT after table creation
-- PostgreSQL would need a separate migration or table recreation for proper FK

-- Add index for efficient querying of uploaded books by user
CREATE INDEX idx_savedbooks_uploaded_by ON savedbooks(uploaded_by);
CREATE INDEX idx_savedbooks_is_uploaded ON savedbooks(is_uploaded);
