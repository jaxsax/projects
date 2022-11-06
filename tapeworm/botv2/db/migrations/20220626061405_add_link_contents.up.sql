CREATE TABLE IF NOT EXISTS link_contents (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    uri VARCHAR(1024) NOT NULL,
    contents TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS link_content_uri_idx ON link_contents(uri);
