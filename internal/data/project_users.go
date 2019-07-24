package data

import (
	"Gomeisa/pkg/utils"
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
	tx, err := utils.Db.Begin()
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

		stmt, err := tx.Prepare(`INSERT into projects_users(project_id, user_id) SELECT $1, users.id FROM users WHERE users.uuid = $2`)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(lastInsertID, userDB.Uuid); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}