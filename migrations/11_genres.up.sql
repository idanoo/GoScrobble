CREATE TABLE IF NOT EXISTS `genres` (
    `uuid` BINARY(16) PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    KEY `nameLookup` (`name`)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;