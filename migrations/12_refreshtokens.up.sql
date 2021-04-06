CREATE TABLE IF NOT EXISTS `refresh_tokens` (
    `uuid` BINARY(16) PRIMARY KEY,
    `user` BINARY(16),
    `token` VARCHAR(64) NOT NULL,
    `expiry` DATETIME NOT NULL DEFAULT NOW(),
    KEY `userLookup` (`user`),
    KEY `tokenLookup` (`token`),
    KEY `expiryLookup` (`expiry`),
    FOREIGN KEY (user) REFERENCES users(uuid)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;