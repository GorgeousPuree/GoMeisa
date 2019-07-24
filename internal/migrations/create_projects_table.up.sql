create table projects
(
    id          serial not null
        constraint projects_pk
            primary key,
    uuid        text   not null,
    name        text   not null,
    description text default 'There is no description yet!'::text
);

alter table projects
    owner to postgres;

create unique index projects_uuid_uindex
    on projects (uuid);