create table specialties
(
    id   serial not null
        constraint specialties_pk
            primary key,
    name text   not null
);

alter table specialties
    owner to postgres;

create unique index specialties_name_uindex
    on specialties (name);