-- +migrate Up

create table notifications (
    id bigserial primary key,
    created_at timestamp without time zone not null default current_timestamp,
    scheduled_for timestamp without time zone not null,
    topic text not null,
    token text,
    priority int not null,
    channel text,
    message jsonb not null
);

create unique index notifications_token_constraint on notifications using btree (token) where token is not null;

create table deliveries (
    id bigserial primary key,
    notification_id int not null references notifications (id) on delete cascade,
    destination text not null,
    destination_type text not null,
    status text not null,
    sent_at timestamp without time zone
);

-- +migrate Down

drop table deliveries;
drop table notifications;

