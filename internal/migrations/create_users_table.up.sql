create table users
(
    id       serial not null
        constraint users_pkey
            primary key,
    uuid     text   not null
        constraint users_uuid_key
            unique,
    email    text   not null
        constraint users_email_key
            unique,
    password text   not null
);

alter table users
    owner to postgres;