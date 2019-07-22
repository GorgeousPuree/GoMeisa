package main

import (
	"Gomeisa"
	"Gomeisa/data"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
	"strconv"
)

// TODO: reduce repetition of code, write common functions.
func projectGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.InProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.GetName()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.GetEmployees()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read employees of project.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var specialty string

	for _, employee := range employees {
		if employee.Email == userDB.Email {
			specialty = employee.Specialty
			break
		}
	}

	type ViewData struct{
		Username string
		Specialty string
		Employees []data.Employee
		Project data.ProjectDB
	}

	data := ViewData{
		Username: userDB.Email,
		Specialty: specialty,
		Project: projectDB,
		Employees: employees,
	}

	Gomeisa.GenerateHTML(w, data, "project_layout", "navbar", "project_main")
}

func inviteGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.InProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.GetName()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.GetEmployees()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read employees of project.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var specialty string

	for _, employee := range employees {
		if employee.Email == userDB.Email {
			specialty = employee.Specialty
			break
		}
	}

	// This operation is provided only for "Technical leader" (admin)
	if specialty != "Technical leader" {
		session.AddFlash("Not enough rights to perform this operation!\n")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/" + projectDB.Uuid + "/", http.StatusSeeOther)
		return
	}

	specialties, err := data.GetSpecialties()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read specialties.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ViewData struct{
		Username string
		Specialty string
		Employees []data.Employee
		Project data.ProjectDB
		Specialties []data.SpecialtyDB
	}

	data := ViewData{
		Username: userDB.Email,
		Specialty: specialty,
		Project: projectDB,
		Employees: employees,
		Specialties: specialties,
	}

	Gomeisa.GenerateHTML(w, data, "project_layout", "navbar", "project_invitation")
}

// TODO: implement invite link output with AJAX
func invitePostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.InProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.GetName()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.GetEmployees()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read employees of project.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var specialty string

	for _, employee := range employees {
		if employee.Email == userDB.Email {
			specialty = employee.Specialty
			break
		}
	}

	// This operation is provided only for "Technical leader" (admin)
	if specialty != "Technical leader" {
		session.AddFlash("Not enough rights to perform this operation!\n")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/" + projectDB.Uuid + "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	specialtyIdStr := r.Form["specialtyId"][0]
	specialtyId, err := strconv.Atoi(specialtyIdStr)

	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to Atoi() number.\n")
		return
	}

	inviteLink, err := projectDB.InviteEmployee(specialtyId)
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to invite employee.\n")
		return
	}

	type ViewData struct{
		Specialty string
		Username string
		Employees []data.Employee
		Project data.ProjectDB
		InviteLink string
	}

	data := ViewData{
		Username: userDB.Email,
		Specialty:specialty,
		Project: projectDB,
		Employees: employees,
		InviteLink:inviteLink,
	}

	Gomeisa.GenerateHTML(w, data, "project_layout", "navbar", "project_inviteLink")
}

func tasksGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.InProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.GetName()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.GetEmployees()
	if err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read employees of project.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	var specialty string

	for _, employee := range employees {
		if employee.Email == userDB.Email {
			specialty = employee.Specialty
			break
		}
	}

	type ViewData struct{
		Username string
		Specialty string
		Employees []data.Employee
		Project data.ProjectDB
	}

	data := ViewData{
		Username: userDB.Email,
		Specialty:specialty,
		Project: projectDB,
		Employees: employees,
	}

	Gomeisa.GenerateHTML(w, data, "project_layout", "navbar", "project_tasks")
}

func joinPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If the user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			Gomeisa.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	if err := userDB.ReadUUID(); err != nil {
		Gomeisa.Error(err, "Error occurred while trying to read user UUID.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	inviteKey := vars["key"]

	projectUUID, err := userDB.Join(inviteKey)

	if err != nil {
		// There is no need to log "duplicate key value violates unique" error
		if _, ok := err.(*pq.Error); ok {
			session.AddFlash("You are already a member of this project!")
			err := session.Save(r, w)

			if err != nil {
				Gomeisa.Error(err, "Error occurred while trying to save session.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return

		} else {
			Gomeisa.Error(err, "Error occurred while trying to join project.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/project/" + projectUUID + "/", http.StatusSeeOther)
}