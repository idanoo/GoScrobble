START TRANSACTION;

CREATE TABLE IF NOT EXISTS `links` (
    `scrobble` BINARY(16) NOT NULL,
    `track` BINARY(16) NOT NULL,
    PRIMARY KEY (`scrobble`, `track`),
    KEY `trackLookup` (`track`)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `artists` (
    `uuid` BINARY(16) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `desc` TEXT,
    `img` VARCHAR(255) DEFAULT NULL
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `albums` (
    `uuid` BINARY(16) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `desc` TEXT,
    `img` VARCHAR(255) DEFAULT NULL
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `tracks` (
    `uuid` BINARY(16) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `desc` TEXT,
    `img` VARCHAR(255) DEFAULT NULL
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `scrobbles` (
    `uuid` BINARY(16) PRIMARY KEY,
    `created_at` DATETIME NOT NULL,
    `created_ip` VARBINARY(16) NULL DEFAULT NULL,
    `user` BINARY(16) NOT NULL,
    `track` BINARY(16) NOT NULL,
    `source` VARCHAR(100) NOT NULL DEFAULT '',
    KEY `userLookup` (`user`),
    KEY `dateLookup` (`created_at`),
    KEY `sourceLookup` (`source`),
    FOREIGN KEY (track) REFERENCES tracks(uuid),
    FOREIGN KEY (user) REFERENCES users(uuid)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `album_artist` (
    `album` BINARY(16) NOT NULL,
    `artist` BINARY(16) NOT NULL,
    PRIMARY KEY (`album`, `artist`),
    FOREIGN KEY (album) REFERENCES albums(uuid),
    FOREIGN KEY (artist) REFERENCES artists(uuid)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `track_album` (
    `track` BINARY(16) NOT NULL,
    `album` BINARY(16) NOT NULL,
    PRIMARY KEY (`track`, `album`),
    FOREIGN KEY (track) REFERENCES tracks(uuid),
    FOREIGN KEY (album) REFERENCES albums(uuid)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `track_artist` (
    `track` BINARY(16) NOT NULL,
    `artist` BINARY(16) NOT NULL,
    PRIMARY KEY (`track`, `artist`),
    FOREIGN KEY (track) REFERENCES tracks(uuid),
    FOREIGN KEY (artist) REFERENCES artists(uuid)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

COMMIT;
