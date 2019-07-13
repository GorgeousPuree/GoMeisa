package main

import (
	"Gomeisa"
	"Gomeisa/data"
	"github.com/gorilla/sessions"
	"net/http"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	userSession := &UserSession {
		Email: email,
		Authenticated: true,
	}

	session.Values["userSession"] = userSession
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/main", http.StatusFound)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "registration.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	// do I need it?
	if err != nil {
		Gomeisa.Danger(err, "Невозможно считать форму!")
	}

	user := data.UserDB{
		Email: r.PostFormValue("email"),
	}

	if err := user.Create(); err != nil {
		Gomeisa.Danger(err, "Невозможно зарегистрировать пользователя")
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	/*rows, err := Gomeisa.Db.Query("SELECT name FROM gomeisa.public.")
	if err != nil {
		return
	}
	//var got []string
	got := make(map[int]string)
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return
		}
		//got = append(got, r)
		got[id] = name
	}
	rows.Close()
	t.Execute(w, got)*/

	templates.ExecuteTemplate(w, "main.html", nil)
}

func logoutPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["userSession"] = UserSession{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userSession := getUserSession(session)
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Войдите в аккаунт!")

		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	r.ParseForm()

	project := data.ProjectDB{Name: r.PostForm.Get("projectName")}
	if err := project.Create(); err != nil {
		Gomeisa.Danger(err, "Невозможно создать проект!")
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




