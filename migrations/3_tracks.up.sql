START TRANSACTION;

CREATE TABLE IF NOT EXISTS links (
    scrobble uuid NOT NULL,
    track uuid NOT NULL,
    PRIMARY KEY (scrobble, track)
);

CREATE INDEX trackLookup ON links (track);

CREATE TABLE IF NOT EXISTS artists (
    uuid uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    "desc" TEXT,
    img VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS albums (
    uuid uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    "desc" TEXT,
    img VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS tracks (
    uuid uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    "desc" TEXT,
    img VARCHAR(255) 
);

CREATE TABLE IF NOT EXISTS scrobbles (
    uuid uuid PRIMARY KEY,
    created_at timestamptz NOT NULL,
    created_ip INET NULL,
    "user" uuid NOT NULL,
    track uuid NOT NULL,
    source VARCHAR(100) NOT NULL DEFAULT '',
    FOREIGN KEY (track) REFERENCES tracks(uuid),
    FOREIGN KEY ("user") REFERENCES users(uuid)
);

CREATE INDEX scrobblesUserLookup ON scrobbles ("user");
CREATE INDEX scrobblesDateLookup ON scrobbles (created_at);
CREATE INDEX scrobblesSourceLookup ON scrobbles (source);

CREATE TABLE IF NOT EXISTS album_artist (
    album uuid NOT NULL,
    artist uuid NOT NULL,
    PRIMARY KEY (album, artist),
    FOREIGN KEY (album) REFERENCES albums(uuid),
    FOREIGN KEY (artist) REFERENCES artists(uuid)
);

CREATE TABLE IF NOT EXISTS track_album (
    track uuid NOT NULL,
    album uuid NOT NULL,
    PRIMARY KEY (track, album),
    FOREIGN KEY (track) REFERENCES tracks(uuid),
    FOREIGN KEY (album) REFERENCES albums(uuid)
);

CREATE TABLE IF NOT EXISTS track_artist (
    track uuid NOT NULL,
    artist uuid NOT NULL,
    PRIMARY KEY (track, artist),
    FOREIGN KEY (track) REFERENCES tracks(uuid),
    FOREIGN KEY (artist) REFERENCES artists(uuid)
);

COMMIT;
