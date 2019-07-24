package api

import (
	"Gomeisa/internal/data"
	"Gomeisa/internal/usession"
	"Gomeisa/pkg/utils"
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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
	userDB := data.UserDB{
		Email: r.PostForm.Get("email"),
	}
	password := r.PostForm.Get("password")

	if !RowExists("SELECT id FROM users WHERE email=$1", userDB.Email) {
		session.AddFlash("Invalid login!")
		err := session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err := userDB.ReadHashedPassword()
	if err != nil {
		utils.Error(err, "Error occurred while trying to save .\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDB.HashedPassword), []byte(password)); err != nil {
		session.AddFlash("Invalid password!")
		err := session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	userSession := &usession.UserSession{
		Email:         userDB.Email,
		Authenticated: true,
	}

	session.Values["userSession"] = userSession
	err = session.Save(r, w)

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
	userDB := data.UserDB{
		Email: r.PostFormValue("email"),
	}

	password := r.PostFormValue("password")
	confirmPassword := r.PostForm.Get("confirmPassword")

	if password != confirmPassword {
		session.AddFlash("Password and confirm password should be same!")
		err := session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	patternEmail := `^\w+@\w+\.\w+$`
	patternPassword := `^.{6,}$`

	if matched, err := regexp.Match(patternEmail, []byte(userDB.Email)); !matched || err != nil {
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

	if matched, err := regexp.Match(patternPassword, []byte(password)); !matched || err != nil {
		session.AddFlash("Password must be at least 6 characters long!")

		err = session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	var err error
	userDB.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 8)
	if err !=nil {
		utils.Error(err, "Error occurred while trying to hash password.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := userDB.Add(); err != nil {
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