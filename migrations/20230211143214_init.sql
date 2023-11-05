-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

create schema blog;

-- BEGIN;
SET lock_timeout = '1min';
SELECT CASE WHEN setting = '0' THEN 'deactivated' ELSE setting || unit END AS lock_timeout FROM pg_settings WHERE name = 'lock_timeout';
-- COMMIT;

create table blog.keyword
(
    id                     bigint                               not null,
    value                  text                                 not null,
    primary key (id) include (value)
);



create unique index order_creation_event_seller_id_srid_timestamp_unx ON rating.order_creation_event (seller_id, srid, timestamp);

grant all on rating.order_creation_event to "api";

create table rating.order_change_event
(
    timestamp              timestamp                            not null,
    srid                   text                                 not null,
    state                  text                                 not null
);

select public.create_hypertable('rating.order_change_event', 'timestamp', chunk_time_interval => INTERVAL '1 month');
create unique index order_change_event_srid_state_timestamp_unx ON rating.order_change_event (srid, state, timestamp);

grant all on rating.order_change_event to "api";

create table rating.feedback
(
    timestamp              timestamp                            not null,
    seller_id              bigint                               not null,
    id                     text                                 not null,
    grade                  smallint                             not null
);

select public.create_hypertable('rating.feedback', 'timestamp', chunk_time_interval => INTERVAL '1 month');
create unique index feedback_seller_id_timestamp_id_unx ON rating.feedback (seller_id, timestamp, id) include (grade);

grant all on rating.feedback to "api";

create table rating.order_creation_event_full
(
    timestamp              timestamp                            not null,
    seller_id              bigint                               not null,
    srid                   text                                 not null,
    nm_id                  bigint                               not null,
    is_defect              bool                                 not null
);

select public.create_hypertable('rating.order_creation_event_full', 'timestamp', chunk_time_interval => INTERVAL '1 month');
create unique index order_creation_event_full_seller_id_srid_timestamp_unx ON rating.order_creation_event_full (seller_id, srid, timestamp) include (is_defect);
create index order_creation_event_full_nm_id_idx ON rating.order_creation_event_full (nm_id) include (seller_id);

grant all on rating.order_creation_event_full to "api";


create table rating.aggregate_for_rating
(
    month                               date                    not null,
    seller_old_id                       bigint                  not null,
    nb_delivered                        bigint                  null,
    nb_orders_marketplace               bigint                  null,
    avg_buyer_rating                    numeric(325, 324)       null,
    nb_buyer_ratings                    bigint                  null,
    nb_defected                         bigint                  null,
    nb_orders_total                     bigint                  null,
    timestamp                           timestamp               not null default now()
);

select public.create_hypertable('rating.aggregate_for_rating', 'month', chunk_time_interval => INTERVAL '1 year');
create unique index aggregate_for_rating__seller_id_month_unx ON rating.aggregate_for_rating (seller_old_id, "month") include (nb_delivered, nb_orders_marketplace, avg_buyer_rating, nb_buyer_ratings, nb_defected, nb_orders_total, timestamp);

grant all on rating.aggregate_for_rating to "api";


create table rating.feedback_processing
(
    seller_id              bigint                               primary key,
    processing_timestamp   timestamp                            null,
    timestamp              timestamp                            not null default now()
);
create index feedback_processing__timestamp_idx ON rating.feedback_processing (timestamp);
create index feedback_processing__processing_timestamp_idx ON rating.feedback_processing (processing_timestamp) include (seller_id, timestamp);

grant all on rating.feedback_processing to "api";

create table rating.defect_processing
(
    srid                   text                                 primary key,
    nm_id                  bigint                               not null,
    timestamp              timestamp                            not null,
    processing_timestamp   timestamp                            null
);

create index defect_processing_processing_timestamp_idx ON rating.defect_processing (processing_timestamp) include (srid, nm_id, timestamp);

grant all on rating.defect_processing to "api";


create table rating.rating_processing
(
    seller_id              bigint                               primary key,
    processing_timestamp   timestamp                            null,
    timestamp              timestamp                            not null default now()
);
create index rating_processing__timestamp_idx ON rating.rating_processing (timestamp) include (seller_id);
create index rating_processing__processing_timestamp_idx ON rating.rating_processing (processing_timestamp) include (seller_id, timestamp);

grant all on rating.rating_processing to "api";


create table rating.aggregate_for_rating_processing
(
    seller_id              bigint                               primary key,
    processing_timestamp   timestamp                            null,
    timestamp              timestamp                            not null default now()
);
create index aggregate_for_rating_processing__timestamp_idx ON rating.aggregate_for_rating_processing (timestamp);
create index aggregate_for_rating_processing__processing_timestamp_idx ON rating.aggregate_for_rating_processing (processing_timestamp) include (seller_id, timestamp);

grant all on rating.aggregate_for_rating_processing to "api";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

drop schema rating;
-- +goose StatementEnd
