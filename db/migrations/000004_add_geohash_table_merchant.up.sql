ALTER TABLE merchants ADD COLUMN geohash VARCHAR(12);

CREATE INDEX merchants_geohash_idx ON merchants (geohash);