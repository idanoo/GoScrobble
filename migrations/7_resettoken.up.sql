CREATE TABLE IF NOT EXISTS `resettoken` (
    `user` BINARY(16) PRIMARY KEY,
    `token` VARCHAR(64) NOT NULL,
    `expiry` DATETIME NOT NULL,
    KEY `tokenLookup` (`token`)
) DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
