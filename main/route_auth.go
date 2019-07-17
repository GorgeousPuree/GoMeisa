package main

import (
	"Gomeisa"
	"Gomeisa/data"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"regexp"
)

func signinGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "signin.html", nil)
}

func signinPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	r.ParseForm()
	email := r.PostForm.Get("email")

	if !Gomeisa.RowExists("SELECT id FROM users WHERE email=$1", email) {
		session.AddFlash("Неверный логин")
		err = session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	userSession := &UserSession{
		Email:         email,
		Authenticated: true,
	}

	session.Values["userSession"] = userSession
	err = session.Save(r, w)

	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/main", http.StatusFound)
}

func signupGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "signup.html", nil)
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

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if _, err := user.Create(); err != nil {
		Gomeisa.Error(err, "User could not be registered/added to database.\n")
		session.AddFlash("User could not be registered!")
		err = session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func mainGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)

	// If the userDB is unauthorized, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in account!")

		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userDB := data.UserDB{Email: userSession.Email}
	projectsUser, err := userDB.GetUserProjects()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read user's projects.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "main.html", projectsUser)
}

func logoutPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["userSession"] = UserSession{}
	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
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

		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

		http.Redirect(w, r, "/main", http.StatusFound)
		return
	}

	project := data.ProjectDB{Name: projectName}
	user := data.UserDB{Email: userSession.Email}

	err := data.CreateProjectUsers(user, project)

	if err != nil {
		Gomeisa.Error(err, "Project could not be added to database.\n")
		log.Println(err)
		session.AddFlash("Project could not be created!")
		session.Save(r, w)
	}

	http.Redirect(w, r, "/main", http.StatusFound)
}

// Need to implement: scanning project uuid from URL.
// Now just passing '1' at 207 line.
func projectGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email:userSession.Email}

	// If the user is unauthorized, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.InProject("1"); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "project.html", nil)
}

func getUserSession(s *sessions.Session) UserSession {
	val := s.Values["userSession"]
	var user = UserSession{}
	user, ok := val.(UserSession)

	if !ok {
		return UserSession{Authenticated: false}
	}

	return user
}