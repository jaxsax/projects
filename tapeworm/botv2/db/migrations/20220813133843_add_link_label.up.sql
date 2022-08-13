CREATE TABLE IF NOT EXISTS link_label (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    link_id INTEGER NOT NULL,
    label_name varchar(64) NOT NULL,
    UNIQUE(link_id, label_name)
);