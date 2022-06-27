CREATE TABLE IF NOT EXISTS refresh_tokens (
    uuid uuid PRIMARY KEY,
    "user" uuid,
    token VARCHAR(64) NOT NULL,
    expiry timestamptz NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("user") REFERENCES users(uuid)
);

CREATE INDEX refreshtokensUserLookup ON refresh_tokens ("user");
CREATE INDEX refreshtokensTokenLookup ON refresh_tokens (token);
CREATE INDEX refreshtokensExpiryLookup ON refresh_tokens (expiry);
