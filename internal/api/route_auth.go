package api

import (
	"Gomeisa/internal/data"
	"Gomeisa/internal/usession"
	"Gomeisa/pkg/utils"
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"log"
	"net/http"
	"regexp"
)

// TODO: exchange method RowExists()
func RowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := utils.Db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}

func SigninGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.GenerateHTML(w, nil, "signin")
}

func SigninPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	r.ParseForm()
	email := r.PostForm.Get("email")

	if !RowExists("SELECT id FROM users WHERE email=$1", email) {
		session.AddFlash("Invalid login!")
		err := session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save .\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	userSession := &usession.UserSession{
		Email:         email,
		Authenticated: true,
	}

	session.Values["userSession"] = userSession
	err := session.Save(r, w)

	if err != nil {
		utils.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects", http.StatusFound)
}

func SignupGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.GenerateHTML(w, nil, "signup")
}

func SignupPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	r.ParseForm()
	user := data.UserDB{
		Email: r.PostFormValue("email"),
	}

	pattern := `^\w+@\w+\.\w+$`

	if matched, err := regexp.Match(pattern, []byte(user.Email)); !matched || err != nil {
		session.AddFlash("Email is not valid!")

		err = session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	if _, err := user.Add(); err != nil {
		// There is no need to log "duplicate key value violates unique" error
		if _, ok := err.(*pq.Error); ok {
			session.AddFlash("Email has already been taken!")
			err := session.Save(r, w)

			if err != nil {
				utils.Error(err, "Error occurred while trying to save session.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return

		} else {
			utils.Error(err, "User could not be registered/added to database.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func LogoutPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")

	session.Values["userSession"] = usession.UserSession{}
	session.Values["projectUUID"] = ""
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		utils.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func getUserSession(s *sessions.Session) usession.UserSession {
	val := s.Values["userSession"]
	var user = usession.UserSession{}
	user, ok := val.(usession.UserSession)

	if !ok {
		return usession.UserSession{Authenticated: false}
	}

	return user
}