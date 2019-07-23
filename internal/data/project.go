package data

import (
	"Gomeisa/pkg/utils"
)

type ProjectDB struct {
	Uuid string
	Name string
	Description string
}

type Employee struct {
	Email      string
	Specialty string
}

type Task struct {
	Description string
	Email string
}

func (projectDB *ProjectDB) ReadName() error {
	row := utils.Db.QueryRow("SELECT name FROM projects WHERE uuid=$1", projectDB.Uuid)
	err := row.Scan(&projectDB.Name)
	return err
}

func (projectDB ProjectDB) ReadEmployees() ([]Employee, error) {
	employees := []Employee{}

	rows, err := utils.Db.Query("SELECT users.email, specialties.name " +
		"FROM users, specialties, projects, projects_users "+
		"WHERE projects.uuid = $1 " +
		"AND projects_users.project_id = projects.id " +
		"AND projects_users.user_id = users.id "+
		"AND projects_users.specialty_id = specialties.id ", projectDB.Uuid)

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

func (projectDB ProjectDB) InviteEmployee(specialtyID int) (string, error) {
	var inviteLink string
	inviteKey := utils.GenerateInviteKey()

	_, err := utils.Db.Query("INSERT into invitations(key, specialty_id, project_id) SELECT $1, $2, id " +
		"FROM projects WHERE uuid = $3;", inviteKey, specialtyID, projectDB.Uuid)
	if err != nil {
		return inviteLink, err
	}

	// TODO: get domain by func instead of writing it manually
	inviteLink = "http://localhost:8080/join/" + inviteKey + "/"
	return inviteLink, err
}

func (projectDB ProjectDB) RemoveEmployee(userDB UserDB) error {
	_, err := utils.Db.Exec("DELETE FROM projects_users " +
		"USING users, projects " +
		"WHERE projects_users.user_id = users.id " +
		"AND users.email = $1 " +
		"AND projects_users.project_id = projects.id " +
		"AND projects.uuid = $2 ",
		userDB.Email, projectDB.Uuid)

	if err != nil {
		return err
	}
	return nil
}

func (projectDB ProjectDB) AddTask(task Task) error {
	_, err := utils.Db.Exec("INSERT into tasks (description, project_id, user_id) " +
		"SELECT $1, p.id, u.id " +
		"FROM projects p, users u " +
		"WHERE p.uuid = $2 " +
		"AND u.email = $3 ",
		task.Description, projectDB.Uuid, task.Email)

	if err != nil {
		return err
	}
	return nil
}

// TODO: implement view for client
func (projectDB ProjectDB) RemoveTask(task Task) error {
	_, err := utils.Db.Exec("DELETE tasks " +
		"WHERE tasks.email = $1, tasks.description = $2",
		task.Email, task.Description)
	
	if err != nil {
		return err 
	}
	return nil
}

func (projectDB ProjectDB) ReadTasks() ([]Task, error) {
	tasks := []Task{}

	rows, err := utils.Db.Query("SELECT users.email, tasks.description " +
		"FROM users, tasks, projects " +
		"WHERE tasks.project_id = projects.id " +
		"AND tasks.user_id = users.id " +
		"AND projects.uuid = $1", projectDB.Uuid)

	if err != nil {
		return tasks, err
	}

	for rows.Next() {
		var task Task

		err = rows.Scan(&task.Email, &task.Description)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	rows.Close()
	return tasks, err
}

func (projectDB ProjectDB) UpdateProjectDescription(projectDescription string ) error {
	_, err := utils.Db.Exec("UPDATE projects " +
		"SET description = $1 " +
		"WHERE uuid = $2",
		projectDescription, projectDB.Uuid)

	if err != nil {
		return err
	}
	return nil
}

func (projectDB *ProjectDB) ReadProjectDescription() error {
	row := utils.Db.QueryRow("SELECT projects.description "+
		"FROM projects "+
		"WHERE projects.uuid = $1",
		projectDB.Uuid)

	err := row.Scan(&projectDB.Description)

	if err != nil {
		return err
	}
	return nil
}