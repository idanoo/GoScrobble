START TRANSACTION;

ALTER TABLE `tracks` DROP COLUMN `spotify_id`;
ALTER TABLE `users` DROP COLUMN `spotify_id`;
ALTER TABLE `albums` DROP COLUMN `spotify_id`;
ALTER TABLE `artists` DROP COLUMN `spotify_id`;

COMMIT;