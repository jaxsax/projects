CREATE TABLE IF NOT EXISTS links_v2 (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    created_ts BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    link VARCHAR(1024) NOT NULL,
    title VARCHAR(1024) NOT NULL,
    extra_data TEXT NOT NULL,
    deleted_at bigint not null default '0',
    host text NOT NULL DEFAULT '', 
    `path` text NOT NULL DEFAULT '', 
    dim_collected INTEGER DEFAULT 0,
    CONSTRAINT uk_link UNIQUE(link)
);
