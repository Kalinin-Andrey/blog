-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
CREATE EXTENSION IF NOT EXISTS btree_gin CASCADE;

create schema blog;

-- BEGIN;
SET lock_timeout = '1min';
SELECT CASE WHEN setting = '0' THEN 'deactivated' ELSE setting || unit END AS lock_timeout FROM pg_settings WHERE name = 'lock_timeout';
-- COMMIT;

create table blog.keyword
(
    id                     integer                              not null,
    sysname                text                                 not null,
    value                  text                                 not null,
    constraint keyword__id__pkey primary key (id) include (sysname, value)
);
create unique index keyword__sysname__unx ON blog.keyword (sysname) include (id, value);

create table blog.tag
(
    id                     integer                              not null,
    sysname                text                                 not null,
    value                  text                                 not null,
    constraint tag__id__pkey primary key (id) include (sysname, value)
);
create unique index tag__sysname__unx ON blog.tag (sysname) include (id, value);

create table blog.blog
(
    id                     integer                              not null,
    sysname                text                                 not null,
    keyword_ids            integer[]                            null,
    tag_ids                integer[]                            null,
    name                   text                                 not null,
    description            text                                 null,
    constraint blog__id__pkey primary key (id) include (sysname, keyword_ids, tag_ids, name, description)
);
create unique index blog__sysname__unx ON blog.blog (sysname) include (id, keyword_ids, tag_ids, name, description);

create table blog.post
(
    id                     bigint                               not null,
    blog_id                integer                              not null,
    sysname                text                                 not null,
    keyword_ids            integer[]                            null,
    tag_ids                integer[]                            null,
    is_deleted             bool                                 not null default false,
    title                  text                                 not null,
    preview                text                                 not null,
    content                text                                 not null,
    content_tsvector       tsvector                             not null,
    created_at             timestamp                            not null,
    updated_at             timestamp                            null,
    deleted_at             timestamp                            null,
    CONSTRAINT post__blog_id__fkey FOREIGN KEY (blog_id) REFERENCES blog.blog(id)
);
create unique index post__id__created_at__pkey ON blog.post (id, created_at) include (blog_id, sysname, keyword_ids, tag_ids, is_deleted, title, preview, content, updated_at, deleted_at);
create unique index post__blog_id__sysname__unx ON blog.post (blog_id, sysname, created_at) include (id, keyword_ids, tag_ids, is_deleted, title, preview, content, updated_at, deleted_at);

create index post__blog_id__keyword_ids__idx ON blog.post (blog_id, keyword_ids) include (id, sysname, tag_ids, is_deleted, title, preview, created_at, updated_at, deleted_at);
create index post__blog_id__tag_ids__idx ON blog.post (blog_id, tag_ids) include (id, sysname, keyword_ids, is_deleted, title, preview, created_at, updated_at, deleted_at);

create index post__blog_id__content_tsvector__ginx ON blog.post using gin (blog_id, content_tsvector, created_at);

select public.create_hypertable('blog.post', 'created_at', chunk_time_interval => INTERVAL '10 year');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop schema blog cascade;
-- +goose StatementEnd
