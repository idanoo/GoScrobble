START TRANSACTION;

ALTER TABLE albums ADD COLUMN `mbid` VARCHAR(36) DEFAULT NULL;
ALTER TABLE artists ADD COLUMN `mbid` VARCHAR(36) DEFAULT NULL;
ALTER TABLE tracks ADD COLUMN `mbid` VARCHAR(36) DEFAULT NULL;

COMMIT;