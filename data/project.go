package data

import "Gomeisa"

type ProjectDB struct {
	Name string
}

func (project *ProjectDB) Create() error {
	_, err := Gomeisa.Db.Exec("INSERT into projects(name) values ($1);", project.Name)
	return err
}

