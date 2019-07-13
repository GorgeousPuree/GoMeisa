package main

import (
	"Gomeisa"
	"database/sql"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var store *sessions.CookieStore
var templates *template.Template

type UserSession struct {
	Email      string
	Authenticated bool
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/logout", logoutPostHandler).Methods("POST")

	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")

	r.HandleFunc("/main", mainGetHandler).Methods("GET")
	r.HandleFunc("/create", createPostHandler).Methods("POST")

	r.HandleFunc("/", mainGetHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func init() {
	rand.Seed(time.Now().UnixNano())
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	// Registering the custom UserSession type with gob encoding package so it can be written as a session value.
	gob.Register(UserSession{})
	templates = template.Must(template.ParseGlob("templates/*.html"))

	var err error
	connStr := "user=postgres password=12345 dbname=gomeisa sslmode=disable"
	Gomeisa.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return
}