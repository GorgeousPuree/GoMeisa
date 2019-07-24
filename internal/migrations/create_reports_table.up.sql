create table reports
(
    id          serial  not null
        constraint reports_pk
            primary key,
    uuid        text    not null,
    project_id  integer not null
        constraint reports_projects_id_fk
            references projects
            on update cascade on delete cascade,
    user_id     integer not null
        constraint reports_users_id_fk
            references users
            on update cascade on delete cascade,
    description integer not null
);

alter table reports
    owner to postgres;

create unique index reports_uuid_uindex
    on reports (uuid);