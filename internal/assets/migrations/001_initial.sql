-- +migrate Up

create table destinations (
    id bigserial primary key,
    type text not null,
    destination text not null
);

create unique index destinations_type_destination on destinations using btree (type,destination);

create table notifications (
    id bigserial primary key,
    created_at timestamp without time zone not null default current_timestamp,
    scheduled_for timestamp without time zone not null,
    topic text not null,
    token text,
    locale text,
    priority int not null,
    delivery_type text,
    message jsonb not null
);

create unique index notifications_token_constraint on notifications using btree (token) where token is not null;

create table delivery_statuses (
    notification_id int not null references notifications (id) on delete cascade,
    destination_id int not null references destinations (id) on delete cascade,
    primary key (notification_id, destination_id),
    status text not null,
    sent_at timestamp without time zone
);



-- +migrate Down

drop table delivery_statuses;
drop table notifications;
drop table destinations;

