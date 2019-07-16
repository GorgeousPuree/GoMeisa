package data

type ProjectDB struct {
	Name string
}

func (project *ProjectDB) Create() (int, error) {
	var lastInsertId int

	rows, err := InsertReturning("INSERT into projects(name) values ($1) RETURNING id", project.Name)
	if err != nil {
		return lastInsertId, err
	}

	err = rows.Scan(&lastInsertId)
	if err != nil {
		return lastInsertId, err
	}

	return lastInsertId, nil
}
