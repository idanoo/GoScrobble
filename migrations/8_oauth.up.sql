CREATE TABLE IF NOT EXISTS `oauth_tokens` (
    `user` BINARY(16),
    `service` VARCHAR(64) NOT NULL,
    `access_token` VARCHAR(255) NULL DEFAULT '',
    `refresh_token` VARCHAR(255) NULL DEFAULT '',
    `expiry` DATETIME NOT NULL,
    `username` VARCHAR(100) NULL DEFAULT '',
    `last_synced` DATETIME NOT NULL DEFAULT NOW(),
    PRIMARY KEY `userService` (`user`, `service`)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
