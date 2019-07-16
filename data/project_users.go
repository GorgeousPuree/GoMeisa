package data

import (
	"Gomeisa"
	"database/sql"
	"log"
)

type ProjectUsersDB struct {
	Uuid        string
	ProjectId   int
	SpecialtyId int
}

// Need to implement: checking on same project names of user
func CreateProjectUsers(userDB UserDB, projectDB ProjectDB) error {
	var lastInsertID int
	// To prevent one SQL-query from executing if another one fails, SQL-transaction was implemented
	tx, err := Gomeisa.Db.Begin()
	if err != nil {
		return err
	}

	{

		/*row, err := InsertReturning("INSERT into projects(name) values ($1) RETURNING id", projectDB.Name)

		err = row.Scan(&lastInsertID)
		if err != nil {
			tx.Rollback()
			return err
		}*/

		// By this way we get unhelpful no rows in result set, if error occurs
		err := tx.QueryRow("INSERT into projects(name) values ($1) RETURNING id", projectDB.Name).Scan(&lastInsertID)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return err
		}

	}

	{
		userUUID, err := userDB.GetUUID()
		if err != nil {
			log.Println(err)
			return err
		}

		stmt, err := tx.Prepare(`INSERT into projects_users(user_uuid, project_id) values ($1, $2);`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(userUUID, lastInsertID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Get slice of strings, which are composed of a project name and specialty of employee
func GetProjectUsers(userDB UserDB) ([]string, error) {
	var got []string
	var err error
	userDB.Uuid, err = userDB.GetUUID()
	if err != nil {
		log.Println(err)
		return got, err
	}

	rows, err := Gomeisa.Db.Query("SELECT projects.name, specialties.name FROM projects, specialties, projects_users, users "+
		"WHERE projects_users.user_uuid = $1 AND projects.id = projects_users.project_id AND specialties.id = projects_users.specialty_id", userDB.Uuid)

	if err != nil {
		return got, err
	}

	for rows.Next() {
		var project string
		var user string
		err = rows.Scan(&project, &user)
		var projectsUser = project + " - " + user
		if err != nil {
			return got, err
		}
		got = append(got, projectsUser)
	}

	rows.Close()
	return got, err
}
