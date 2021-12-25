CREATE TABLE IF NOT EXISTS oauth_tokens (
    "user" uuid,
    service VARCHAR(64) NOT NULL,
    access_token VARCHAR(255) NULL DEFAULT '',
    refresh_token VARCHAR(255) NULL DEFAULT '',
    url VARCHAR(255) NULL DEFAULT '',
    expiry timestamptz NOT NULL,
    username VARCHAR(100) NULL DEFAULT '',
    last_synced timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("user", "service")
);
