package data

import (
	"Gomeisa"
)

type ProjectDB struct {
	Uuid string
	Name string
}

type Employee struct {
	Email      string
	Specialty string
}

func (projectDB *ProjectDB) GetName() error {
	row := Gomeisa.Db.QueryRow("SELECT name FROM projects WHERE uuid=$1", projectDB.Uuid)
	err := row.Scan(&projectDB.Name)
	return err
}

func (projectDB *ProjectDB) GetEmployees() ([]Employee, error) {
	employees := []Employee{}

	rows, err := Gomeisa.Db.Query("SELECT users.email, specialties.name FROM users, specialties, projects, projects_users "+
		"WHERE projects.uuid = $1 AND projects_users.project_id = projects.id AND projects_users.user_uuid = users.uuid "+
		"AND projects_users.specialty_id = specialties.id", projectDB.Uuid)

	if err != nil {
		return employees, err
	}

	for rows.Next() {
		var employee Employee

		err = rows.Scan(&employee.Email, &employee.Specialty)
		if err != nil {
			return employees, err
		}

		employees = append(employees, employee)
	}
	rows.Close()
	return employees, err
}

func (projectDB *ProjectDB) InviteEmployee(specialtyID int) (string, error) {
	var inviteLink string
	inviteKey := Gomeisa.GenerateInviteKey()

	_, err := Gomeisa.Db.Query("INSERT into invitations(key, specialty_id, project_id) SELECT $1, $2, id " +
		"FROM projects WHERE uuid = $3;", inviteKey, specialtyID, projectDB.Uuid)
	if err != nil {
		return inviteLink, err
	}

	// TODO: get domain by func instead of writing it manually
	inviteLink = "http://localhost:8080/join/" + inviteKey + "/"
	return inviteLink, err
}