-- Remove redemption codes system
-- Compatible with both SQLite and PostgreSQL

-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS redemption_log;
DROP TABLE IF EXISTS redemption_codes;
