CREATE TABLE IF NOT EXISTS users (
    uuid uuid PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    created_ip INET NULL,
    modified_at timestamptz NOT NULL DEFAULT NOW(),
    modified_ip INET NULL,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(60) NOT NULL,
    email VARCHAR(255) NULL,
    verified BOOL NOT NULL DEFAULT FALSE,
    active BOOL NOT NULL DEFAULT TRUE,
    admin BOOL NOT NULL DEFAULT FALSE,
    private BOOL NOT NULL DEFAULT FALSE,
    timezone VARCHAR(100) NOT NULL DEFAULT 'Pacific/Auckland'
);

CREATE INDEX usersUsernameLookup ON users (username, active);
CREATE INDEX usersEmailLookup ON users (email, active);