package Gomeisa

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
)

//var Templates *template.Template
var Db *sql.DB
var logger *log.Logger
//var symbols = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type UserSession struct {
	Email         string
	Authenticated bool
}

func init() {
	file, err := os.OpenFile("gomeisa.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(args ...interface {}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

/*
// used for generating cookie id and invite key
func GenerateString(length int) string {

	b := make([]rune, length)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}*/

func RowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := Db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists '%s' %v", args, err)
	}
	return exists
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
		files = append(files, fmt.Sprintf("Templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	err := templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		Error(err)
		return
	}
}
