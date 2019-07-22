package main

import (
	"Gomeisa"
	"database/sql"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var store *sessions.CookieStore

func main() {
	r := mux.NewRouter()
	projectRouter := r.PathPrefix("/project/{key:[^ ]+}/").Subrouter()
	inviteRouter := projectRouter.PathPrefix("/invite").Subrouter()
	taskRouter := projectRouter.PathPrefix("/tasks").Subrouter()

	r.HandleFunc("/signin", signinGetHandler).Methods("GET")
	r.HandleFunc("/signin", signinPostHandler).Methods("POST")
	r.HandleFunc("/logout", logoutPostHandler).Methods("POST")

	r.HandleFunc("/signup", signupGetHandler).Methods("GET")
	r.HandleFunc("/signup", signupPostHandler).Methods("POST")
	r.HandleFunc("/join/{key:[^ ]+}/", joinPostHandler)

	r.HandleFunc("/projects", projectsGetHandler).Methods("GET")
	r.HandleFunc("/createProject", createProjectPostHandler).Methods("POST")

	projectRouter.HandleFunc("", projectGetHandler).Methods("GET")
	inviteRouter.HandleFunc("", inviteGetHandler).Methods("GET")
	inviteRouter.HandleFunc("", invitePostHandler).Methods("POST")
	taskRouter.HandleFunc("", tasksGetHandler).Methods("GET")


	r.HandleFunc("/", projectsGetHandler)
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
	gob.Register(Gomeisa.UserSession{})

	var err error
	connStr := "user=postgres password=12345 dbname=gomeisa sslmode=disable"
	Gomeisa.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to open database.\n")
		log.Fatal(err)
	}
	return
}
