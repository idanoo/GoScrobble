START TRANSACTION;

ALTER TABLE `users` ADD COLUMN `spotify_id` VARCHAR(255) DEFAULT '';
ALTER TABLE `albums` ADD COLUMN `spotify_id` VARCHAR(255) DEFAULT '';
ALTER TABLE `artists` ADD COLUMN `spotify_id` VARCHAR(255) DEFAULT '';
ALTER TABLE `tracks` ADD COLUMN `spotify_id` VARCHAR(255) DEFAULT '';

ALTER TABLE `users` ADD INDEX `spotifyLookup` (`spotify_id`);
ALTER TABLE `albums` ADD INDEX `spotifyLookup` (`spotify_id`);
ALTER TABLE `artists` ADD INDEX `spotifyLookup` (`spotify_id`);
ALTER TABLE `tracks` ADD INDEX `spotifyLookup` (`spotify_id`);

COMMIT;