-- Reverse migration: Remove the geohash column and index from the merchants table

-- Drop the index on the geohash column
DROP INDEX IF EXISTS merchants_geohash_idx;

-- Drop the geohash column
ALTER TABLE merchants DROP COLUMN geohash;

