CREATE TABLE IF NOT EXISTS content_blocklist (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    strategy varchar(64) NOT NULL,
    content varchar(256) NOT NULL
);
