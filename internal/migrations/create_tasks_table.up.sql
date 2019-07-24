create table tasks
(
    id          serial  not null
        constraint tasks_pk
            primary key,
    user_id     integer not null
        constraint tasks_users_id_fk
            references users
            on update cascade on delete cascade,
    project_id  integer not null
        constraint tasks_projects_id_fk
            references projects
            on update cascade on delete cascade,
    description text    not null
);

alter table tasks
    owner to postgres;
