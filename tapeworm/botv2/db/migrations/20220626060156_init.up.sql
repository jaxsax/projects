CREATE TABLE IF NOT EXISTS links (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    created_ts BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    link VARCHAR(1024) NOT NULL,
    title VARCHAR(1024) NOT NULL,
    extra_data TEXT NOT NULL,
    deleted_at bigint not null default '0'
);

CREATE TABLE IF NOT EXISTS skipped_links (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    error_reason VARCHAR(512) NOT NULL,
    link VARCHAR(1024) NOT NULL
);

CREATE TABLE IF NOT EXISTS updates (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
   "data" TEXT NOT NULL
);
