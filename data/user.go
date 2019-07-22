package data

import (
	"Gomeisa"
	"github.com/nu7hatch/gouuid"
)

type UserDB struct {
	Uuid  string
	Email string
	//Name      string
	//Password  string
	//CreatedAt time.Time
}

func (userDB *UserDB) Create() (int, error) {
	var lastInsertId int
	uuidBytes, err := uuid.NewV4()
	if err != nil {
		return lastInsertId, err
	}

	userDB.Uuid = uuidBytes.String()

	err = Gomeisa.Db.QueryRow("INSERT into users(uuid, email) values ($1, $2) RETURNING id", userDB.Uuid, userDB.Email).Scan(&lastInsertId)
	if err != nil {
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (userDB *UserDB) InProject(projectUUID string) bool {
	var exists bool
	err := userDB.ReadUUID()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read data from db.\n")
	}

	err = Gomeisa.Db.QueryRow("SELECT EXISTS(SELECT pu.project_id FROM projects_users as pu "+
		"JOIN projects as p ON p.id = pu.project_id "+
		"AND pu.user_uuid = $1 AND p.uuid = $2)", userDB.Uuid, projectUUID).Scan(&exists)

	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read data from db.\n")
		return exists
	}
	return exists
}

// Reading UUID of user from db, writing it to userDB.Uuid and returning error.
// Use userDB.Uuid to get an access to filled with ReadUUID() method uuid.
func (userDB *UserDB) ReadUUID() error {
	row := Gomeisa.Db.QueryRow("SELECT uuid FROM users WHERE email=$1", userDB.Email)
	err := row.Scan(&userDB.Uuid)
	return err
}

// Get map with project uuid as key and project name with specialty name of user as value
func (userDB *UserDB) GetUserProjects() (map[string]string, error) {
	got := make(map[string]string)
	err := userDB.ReadUUID()
	if err != nil {
		return got, err
	}

	// Selecting all userDB's projects
	rows, err := Gomeisa.Db.Query("SELECT projects.uuid, projects.name, specialties.name FROM projects, projects_users, specialties "+
		"WHERE projects_users.user_uuid = $1 AND projects_users.project_id = projects.id "+
		"AND projects_users.specialty_id = specialties.id", userDB.Uuid)

	if err != nil {
		return got, err
	}

	for rows.Next() {
		var projectUUID string
		var projectName string
		var userSpecialty string
		err = rows.Scan(&projectUUID, &projectName, &userSpecialty)
		userProjectsAndSpecialties := projectName + " - " + userSpecialty
		if err != nil {
			return got, err
		}
		got[projectUUID] = userProjectsAndSpecialties
	}

	rows.Close()
	return got, err
}

func (userDB *UserDB) Join(inviteKey string) (string, error) {
	var projectId int
	var projectUUID string

	tx, err := Gomeisa.Db.Begin()
	if err != nil {
		return projectUUID, err
	}

	{
		err := tx.QueryRow("INSERT into projects_users (user_uuid, project_id, specialty_id)"+
			"SELECT $1, project_id, specialty_id FROM invitations WHERE key = $2 returning project_id", userDB.Uuid, inviteKey).Scan(&projectId)

		if err != nil  {
			tx.Rollback()
			return projectUUID, err
		}
	}

	{
		err := tx.QueryRow("SELECT uuid FROM projects WHERE id = $1", projectId).Scan(&projectUUID)

		if err != nil {
			tx.Rollback()
			return projectUUID, err
		}
	}

	{
		_, err := tx.Exec("DELETE FROM invitations WHERE key = $1", inviteKey)

		if err != nil {
			tx.Rollback()
			return projectUUID, err
		}
	}
	return projectUUID, tx.Commit()
}