create table projects_users
(
    id           serial            not null
        constraint projects_users_pk
            primary key,
    user_id      integer           not null
        constraint projects_users_users_id_fk
            references users
            on update cascade on delete cascade,
    project_id   integer           not null
        constraint projects_users_projects_id_fk
            references projects
            on update cascade on delete cascade,
    specialty_id integer default 4 not null
        constraint projects_users_specialties_id_fk
            references specialties
            on update cascade on delete cascade,
    constraint projects_users_user_id_project_id_key
        unique (user_id, project_id)
);

alter table projects_users
    owner to postgres;

