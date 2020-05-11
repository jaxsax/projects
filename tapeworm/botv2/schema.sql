CREATE TABLE links (
    id serial primary key,
    created_ts bigint NOT NULL,
    created_by bigint NOT NULL,
    link character varying(1024) NOT NULL,
    title character varying(1024) NOT NULL,
    extra_data jsonb NOT NULL
);

CREATE TABLE skipped_links (
    id serial primary key,
    error_reason character varying(512) NOT NULL,
    link character varying(1024) NOT NULL
);

CREATE TABLE updates (
    id serial primary key,
    "data" jsonb NOT NULL
);
