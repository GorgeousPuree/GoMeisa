package Gomeisa

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var Db *sql.DB
var logger *log.Logger
//var symbols = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
