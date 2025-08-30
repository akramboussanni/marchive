-- Drop all tables and indexes in reverse order of dependencies

-- Drop tables with foreign key dependencies first
DROP TABLE IF EXISTS downloadjobs;
DROP TABLE IF EXISTS savedbooks;
DROP TABLE IF EXISTS downloadrequests;

-- Drop security-related tables
DROP TABLE IF EXISTS lockouts;
DROP TABLE IF EXISTS failed_logins;
DROP TABLE IF EXISTS jwt_blacklist;

-- Drop main users table last
DROP TABLE IF EXISTS users;

