START TRANSACTION;

ALTER TABLE users ADD COLUMN spotify_id VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE albums ADD COLUMN spotify_id VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE artists ADD COLUMN spotify_id VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE tracks ADD COLUMN spotify_id VARCHAR(255) NOT NULL DEFAULT '';

CREATE INDEX usersSpotifyLookup ON users (spotify_id);
CREATE INDEX albumsSpotifyLookup ON albums (spotify_id);
CREATE INDEX artistsSpotifyLookup ON artists (spotify_id);
CREATE INDEX tracksSpotifyLookup ON tracks (spotify_id);

COMMIT;