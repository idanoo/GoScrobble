CREATE TABLE IF NOT EXISTS genres (
    uuid uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE INDEX genresNameLookup ON genres (name)