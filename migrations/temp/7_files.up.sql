CREATE TABLE IF NOT EXISTS `files` (
    `uuid` BINARY(16) PRIMARY KEY,
    `path` VARCHAR(255) NOT NULL,
    `filesize` INT NULL DEFAULT NULL,
    `dimension` VARCHAR(4) NOT NULL,
    KEY `dimensionLookup` (`uuid`, `dimension`)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;