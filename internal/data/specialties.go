package data

import (
	"Gomeisa/pkg/utils"
)

type SpecialtyDB struct {
	Id  string
	Name string
	//Name      string
	//Password  string
	//CreatedAt time.Time
}

func GetSpecialties() ([]SpecialtyDB, error) {
	got := []SpecialtyDB{}

	// Selecting all userDB's projects
	rows, err := utils.Db.Query("SELECT * FROM specialties WHERE name != 'Технический лидер'")

	if err != nil {
		return got, err
	}

	for rows.Next() {
		var specialty SpecialtyDB
		err = rows.Scan(&specialty.Id, &specialty.Name)
		if err != nil {
			return got, err
		}
		got = append(got, specialty)
	}
	rows.Close()
	return got, err
}