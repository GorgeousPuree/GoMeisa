package utils

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var Db *sql.DB
var Logger *log.Logger

func Error(args ...interface {}) {
	Logger.SetPrefix("ERROR ")
	Logger.Println(args...)
}

func GenerateInviteKey() string {
	var numbers = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := 16
	b := make([]rune, length)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(b)
}


func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("../web/templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		Error(err)
		return
	}
}

func ConnectDB() error {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT_BD"),
		os.Getenv("SERVICE_USER"), os.Getenv("SERVICE_PASSWORD"), os.Getenv("SERVICE_DBNAME"))

	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	err = Db.Ping()
	if err != nil {
		return err
	}
	return nil
}