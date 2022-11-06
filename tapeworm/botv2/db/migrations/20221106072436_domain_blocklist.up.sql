CREATE TABLE IF NOT EXISTS domain_blocklist (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    domain varchar(512) NOT NULL,
    UNIQUE(domain)
);
