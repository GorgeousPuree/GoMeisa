package migrations

import (
	"Gomeisa/pkg/utils"
	"database/sql"
	"github.com/lopezator/migrator"
	"log"
)

func Up() {
	m := migrator.New(
		&migrator.MigrationNoTx{
			Name: "Create table users",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table IF NOT EXISTS users 
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
				owner to postgres;`); err != nil {
					return err
				}
				return nil
			},
		},

		&migrator.MigrationNoTx{
			Name: "Create table projects",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table IF NOT EXISTS projects
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
    on projects (uuid);`); err != nil {
					return err
				}
				return nil
			},
		},

		&migrator.MigrationNoTx{
			Name: "Create table specialties",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table IF NOT EXISTS specialties
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
insert into specialties(name) values ('Programmer');
insert into specialties(name) values ('Tester');
insert into specialties(name) values ('Manager');
insert into specialties(name) values ('Technical leader');`); err != nil {
					return err
				}
				return nil
			},
		},

		&migrator.MigrationNoTx{
			Name: "Create table tasks",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table tasks
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
`); err != nil {
					return err
				}
				return nil
			},
		},

		&migrator.MigrationNoTx{
			Name: "Create table projects_users",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table IF NOT EXISTS projects_users
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
`); err != nil {
					return err
				}
				return nil
			},
		},

		&migrator.MigrationNoTx{
			Name: "Create table reports",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table IF NOT EXISTS reports
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
    on reports (uuid);`); err != nil {
					return err
				}
				return nil
			},
		},

		&migrator.MigrationNoTx{
			Name: "Create table invitations",
			Func: func(db *sql.DB) error {
				if _, err := db.Exec(`create table IF NOT EXISTS invitations
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
`); err != nil {
					return err
				}
				return nil
			},
		},
	)

	if err := m.Migrate(utils.Db); err != nil {
		log.Fatal(err)
	}
}
