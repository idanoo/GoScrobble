CREATE TABLE IF NOT EXISTS resettoken (
    "user" uuid PRIMARY KEY,
    token VARCHAR(64) NOT NULL,
    expiry timestamptz NOT NULL
);

CREATE INDEX resettokenTokenLookup ON resettoken (token)
