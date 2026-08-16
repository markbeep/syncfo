-- +goose up
CREATE TABLE threads (
    id sqlite3_int64 PRIMARY KEY,
    guild_id sqlite3_int64 NOT NULL
);

CREATE TABLE messages (
    id sqlite3_int64 PRIMARY KEY,
    thread_id sqlite3_int64 NOT NULL REFERENCES threads(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE messages;
DROP TABLE threads;
