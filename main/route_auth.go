package main

import (
	"Gomeisa"
	"Gomeisa/data"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"regexp"
)

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	r.ParseForm()
	email := r.PostForm.Get("email")

	if !Gomeisa.RowExists("SELECT id FROM users WHERE email=$1", email) {
		session.AddFlash("Неверный логин")
		err = session.Save(r, w)
		if err != nil {
			Gomeisa.Danger(err, "Ошибка сохранения сессии.")
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
		Gomeisa.Danger(err, "Ошибка сохранения сессии.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/main", http.StatusFound)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "registration.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	r.ParseForm()
	user := data.UserDB{
		Email: r.PostFormValue("email"),
	}

	pattern := `^\w+@\w+\.\w+$`

	if matched, err := regexp.Match(pattern, []byte(user.Email)); !matched || err != nil {
		session.AddFlash("Недопустимый email!")

		err = session.Save(r, w)
		if err != nil {
			Gomeisa.Danger(err, "Ошибка сохранения сессии.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}

	if _, err := user.Create(); err != nil {
		Gomeisa.Danger(err, "Невозможно зарегистрировать пользователя!")
		session.AddFlash("Невозможно зарегистрировать пользователя!")
		err = session.Save(r, w)
		if err != nil {
			Gomeisa.Danger(err, "Ошибка сохранения сессии.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func mainGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)

	// If the user is unauthorized, add flash message and return an unauthorized status
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Войдите в аккаунт!")

		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Danger(err, "Ошибка сохранения сессии.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user := data.UserDB{Email: userSession.Email}
	projectsUser, err := data.GetProjectUsers(user)
	if err != nil {
		Gomeisa.Danger(err, "Ошибка считывания проектов пользователя")
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
		Gomeisa.Danger(err, "Ошибка сохранения сессии.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	userSession := getUserSession(session)

	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Войдите в аккаунт!")

		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Danger(err, "Ошибка сохранения сессии.")
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
		log.Printf("Недопустимое имя проекта! %s\n", err)
		session.AddFlash("Недопустимое имя проекта!")

		err := session.Save(r, w)
		if err != nil {
			Gomeisa.Danger(err, "Ошибка сохранения сессии.")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/main", http.StatusFound)
		return
	}

	project := data.ProjectDB{Name: projectName}
	user := data.UserDB{Email: userSession.Email}

	err = data.CreateProjectUsers(user, project)

	if err != nil {
		Gomeisa.Danger(err, "Невозможно создать проект!")
		log.Println(err)
		session.AddFlash("Невозможно создать проект!")
		session.Save(r, w)
	}

	http.Redirect(w, r, "/main", http.StatusFound)
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
