package data

type ProjectDB struct {
	Uuid string
	Name string
}

/*
func (projectDB *ProjectDB) Create() (int, error) {
	var lastInsertId int
	uuidBytes, err := uuid.NewV4()

	if err != nil {
		return lastInsertId, err
	}

	projectDB.Uuid = uuidBytes.String()

	rows, err := InsertReturning("INSERT into projects(name) values ($1) RETURNING id", projectDB.Name)
	if err != nil {
		return lastInsertId, err
	}

	err = rows.Scan(&lastInsertId)
	if err != nil {
		return lastInsertId, err
	}

	err = Gomeisa.Db.QueryRow("INSERT into projects(name) values ($1) RETURNING id", projectDB.Name).Scan(&lastInsertId)
	if err != nil {
		return lastInsertId, err
	}
	return lastInsertId, nil
}
*/

/*func (projectDB *ProjectDB) GetUUID() error {
	row := Gomeisa.Db.QueryRow("SELECT uuid FROM projects WHERE email=$1", userDB.Email)
	err := row.Scan(userDB.Uuid)
	return err
}*/