-- Remove ghost mode columns from savedbooks
DROP INDEX IF EXISTS idx_savedbooks_requested_by;
DROP INDEX IF EXISTS idx_savedbooks_is_ghost;

ALTER TABLE savedbooks DROP COLUMN IF EXISTS requested_by;
ALTER TABLE savedbooks DROP COLUMN IF EXISTS is_ghost;
