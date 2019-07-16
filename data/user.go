package data

import (
	"Gomeisa"
)

type UserDB struct {
	Uuid string
	Email string
	//Name      string
	//Password  string
	//CreatedAt time.Time
}

func (user *UserDB) Create() (int, error) {
	var lastInsertId int
	uuid := Gomeisa.GenerateString(32)

	rows, err := InsertReturning("INSERT into users(uuid, email) values ($1, $2) RETURNING id", uuid, user.Email)
	if err != nil {
		return lastInsertId, err
	}

	err = rows.Scan(&lastInsertId)
	if err != nil {
		return lastInsertId, err
	}

	return lastInsertId, nil
}

func (user *UserDB) GetUUID() (string, error) {
	var userUUID string
	row := Gomeisa.Db.QueryRow("SELECT uuid FROM Users WHERE email=$1", user.Email)
	err := row.Scan(&userUUID)
	return userUUID, err
}