package data

import (
	"Gomeisa"
)

type UserDB struct {
	Id   int
	Uuid string
	//Name      string
	Email string
	//Password  string
	//CreatedAt time.Time
}

func (user *UserDB) Create() error{
	uuid := Gomeisa.GenerateString(32)
	_, err := Gomeisa.Db.Exec("INSERT into users(uuid, email) values ($1, $2);", uuid, user.Email)
	return err
}


