CREATE TABLE IF NOT EXISTS `users` (
    `uuid` BINARY(16) PRIMARY KEY,
    `created_at` DATETIME NOT NULL DEFAULT NOW(),
    `created_ip` VARBINARY(16) NULL DEFAULT NULL,
    `modified_at` DATETIME NOT NULL DEFAULT NOW(),
    `modified_ip` VARBINARY(16) NULL DEFAULT NULL,
    `username` VARCHAR(64) NOT NULL,
    `password` VARCHAR(60) NOT NULL,
    `email` VARCHAR(255) NULL DEFAULT NULL,
    `verified` TINYINT(1) NOT NULL DEFAULT 0,
    `active` TINYINT(1) NOT NULL DEFAULT 1,
    `admin` TINYINT(1) NOT NULL DEFAULT 0,
    `private` TINYINT(1) NOT NULL DEFAULT 0,
    `timezone` VARCHAR(100) NOT NULL DEFAULT 'Etc/UTC',
    KEY `usernameLookup` (`username`, `active`),
    KEY `emailLookup` (`email`, `active`)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;