-- Rollback upload support from savedbooks table

DROP INDEX IF EXISTS idx_savedbooks_is_uploaded;
DROP INDEX IF EXISTS idx_savedbooks_uploaded_by;

ALTER TABLE savedbooks DROP CONSTRAINT IF EXISTS fk_savedbooks_uploaded_by;

ALTER TABLE savedbooks DROP COLUMN IF EXISTS original_filename;
ALTER TABLE savedbooks DROP COLUMN IF EXISTS uploaded_by;
ALTER TABLE savedbooks DROP COLUMN IF EXISTS is_uploaded;
