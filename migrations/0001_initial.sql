-- +goose up
CREATE TABLE threads (
    id sqlite3_int64 PRIMARY KEY,
    guild_id sqlite3_int64 NOT NULL
);

CREATE TABLE messages (
    id sqlite3_int64 PRIMARY KEY,
    thread_id sqlite3_int64 NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    built BOOLEAN NOT NULL DEFAULT FALSE,
    last_edited TIMESTAMP
);

CREATE TABLE tracked_channels (
    id sqlite3_int64 PRIMARY KEY,
    guild_id sqlite3_int64 NOT NULL
);

-- +goose down
DROP TABLE messages;
DROP TABLE threads;

DROP TABLE tracked_channels;
