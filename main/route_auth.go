package main

import (
	"Gomeisa"
	"Gomeisa/data"
	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	"net/http"
	"regexp"
)

func signinGetHandler(w http.ResponseWriter, r *http.Request) {
	Gomeisa.GenerateHTML(w, nil, "signin")
}

func signinPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	r.ParseForm()
	email := r.PostForm.Get("email")

	if !Gomeisa.RowExists("SELECT id FROM users WHERE email=$1", email) {
		session.AddFlash("Invalid login!")
		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	userSession := &Gomeisa.UserSession{
		Email:         email,
		Authenticated: true,
	}

	session.Values["userSession"] = userSession
	err := session.Save(r, w)

	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects", http.StatusFound)
}

func signupGetHandler(w http.ResponseWriter, r *http.Request) {
	Gomeisa.GenerateHTML(w, nil, "signup")
}

func signupPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	r.ParseForm()
	user := data.UserDB{
		Email: r.PostFormValue("email"),
	}

	pattern := `^\w+@\w+\.\w+$`

	if matched, err := regexp.Match(pattern, []byte(user.Email)); !matched || err != nil {
		session.AddFlash("Email is not valid!")

		err = session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	if _, err := user.Create(); err != nil {
		// There is no need to log "duplicate key value violates unique" error
		if _, ok := err.(*pq.Error); ok {
			session.AddFlash("Email has already been taken!")
			err := session.Save(r, w)

			if err != nil {
				Gomeisa.Error(err, "Error occurred while trying to save session.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return

		} else {
			Gomeisa.Error(err, "User could not be registered/added to database.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func projectsGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)

	err := session.Save(r, w)
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the userDB is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in account!")

		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	userDB := data.UserDB{Email: userSession.Email}
	userProjects, err := userDB.GetUserProjects()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read user's projects.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ViewData struct {
		Username string
		UserProjects map[string]string
	}

	data := ViewData{
		Username: userDB.Email,
		UserProjects: userProjects,
	}

	Gomeisa.GenerateHTML(w, data, "projects_layout", "navbar")
}

func logoutPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["userSession"] = Gomeisa.UserSession{}
	session.Values["projectUUID"] = ""
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func createProjectPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)

	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	projectName := r.PostForm.Get("projectName")
	pattern := `.*\S.*`

	if matched, err := regexp.Match(pattern, []byte(projectName)); !matched || err != nil {
		session.AddFlash("Project name is not valid!")

		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/projects", http.StatusFound)
		return
	}

	project := data.ProjectDB{Name: projectName}
	user := data.UserDB{Email: userSession.Email}

	err := data.CreateProjectUsers(user, project)

	if err != nil {
		Gomeisa.Error(err, "Project could not be added to database.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects", http.StatusFound)
}

func getUserSession(s *sessions.Session) Gomeisa.UserSession {
	val := s.Values["userSession"]
	var user = Gomeisa.UserSession{}
	user, ok := val.(Gomeisa.UserSession)

	if !ok {
		return Gomeisa.UserSession{Authenticated: false}
	}

	return user
}