package data

import (
	"Gomeisa"
	uuid "github.com/nu7hatch/gouuid"
)

type ProjectUsersDB struct {
	Uuid        string
	ProjectId   int
	SpecialtyId int
}

// TODO: checking whether user tries to create a project with name which he has already used.
// TODO: regenerating project UUID if it duplicates
func CreateProjectUsers(userDB UserDB, projectDB ProjectDB) error {
	var lastInsertID int

	uuidBytes, err := uuid.NewV4()
	if err != nil {
		return err
	}

	projectDB.Uuid = uuidBytes.String()
	// To prevent one SQL-query from executing if another one fails, SQL-transaction was implemented
	tx, err := Gomeisa.Db.Begin()
	if err != nil {
		return err
	}

	{
		err := tx.QueryRow("INSERT into projects(uuid, name) values ($1, $2) RETURNING id", projectDB.Uuid, projectDB.Name).Scan(&lastInsertID)
		if err != nil  {
			tx.Rollback()
			return err
		}

	}

	{
		err := userDB.ReadUUID()
		if err != nil {
			tx.Rollback()
			return err
		}

		stmt, err := tx.Prepare(`INSERT into projects_users(user_uuid, project_id) values ($1, $2);`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(userDB.Uuid, lastInsertID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}