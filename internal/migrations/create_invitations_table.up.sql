create table invitations
(
    id           serial  not null
        constraint invitations_pk
            primary key,
    key          text    not null,
    project_id   integer not null
        constraint invitations_project___fk
            references projects
            on update cascade on delete cascade,
    specialty_id integer not null
        constraint invitations_specialty___fk
            references specialties
            on update cascade on delete cascade
);

alter table invitations
    owner to postgres;

create unique index invitations_key_uindex
    on invitations (key);
