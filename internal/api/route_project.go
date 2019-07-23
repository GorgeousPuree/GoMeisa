package api

import (
	"Gomeisa/internal/data"
	"Gomeisa/internal/usession"
	"Gomeisa/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"net/http"
	"regexp"
	"strconv"
)

func ProjectsGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)

	err := session.Save(r, w)
	if err != nil {
		utils.Error(err, "Error occurred while trying to save session.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the userDB is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in account!")

		err := session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	userDB := data.UserDB{Email: userSession.Email}
	userProjects, err := userDB.GetUserProjects()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read user's projects.\n")
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

	utils.GenerateHTML(w, data, "projects_layout", "navbar")
}

func ProjectPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)

	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
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
			utils.Error(err, "Error occurred while trying to save session.\n")
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
		utils.Error(err, "Project could not be added to database.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects", http.StatusFound)
}

// TODO: reduce repetition of code, write common functions.
func ProjectGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}

	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project's name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = projectDB.ReadProjectDescription()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project's description.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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

	type ViewData struct {
		Username  string
		Specialty string
		Employees []data.Employee
		Project   data.ProjectDB
	}

	data := ViewData{
		Username:  userDB.Email,
		Specialty: specialty,
		Project:   projectDB,
		Employees: employees,
	}

	utils.GenerateHTML(w, data, "project_layout", "navbar", "project_main")
}

func InviteGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
		return
	}

	specialties, err := data.GetSpecialties()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read specialties.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ViewData struct {
		Username    string
		Specialty   string
		Employees   []data.Employee
		Project     data.ProjectDB
		Specialties []data.SpecialtyDB
	}

	data := ViewData{
		Username:    userDB.Email,
		Specialty:   specialty,
		Project:     projectDB,
		Employees:   employees,
		Specialties: specialties,
	}

	utils.GenerateHTML(w, data, "project_layout", "navbar", "project_invitation")
}

// TODO: implement invite link output with AJAX
func InvitePostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
		return
	}

	r.ParseForm()

	specialtyIdStr := r.Form["specialtyId"][0]
	specialtyId, err := strconv.Atoi(specialtyIdStr)

	if err != nil {
		utils.Error(err, "Error occurred while trying to Atoi() number.\n")
		return
	}

	inviteLink, err := projectDB.InviteEmployee(specialtyId)
	if err != nil {
		utils.Error(err, "Error occurred while trying to invite employee.\n")
		return
	}

	type ViewData struct {
		Specialty  string
		Username   string
		Employees  []data.Employee
		Project    data.ProjectDB
		InviteLink string
	}

	data := ViewData{
		Username:   userDB.Email,
		Specialty:  specialty,
		Project:    projectDB,
		Employees:  employees,
		InviteLink: inviteLink,
	}

	utils.GenerateHTML(w, data, "project_layout", "navbar", "project_inviteLink")
}

func RemoveEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]
	email := vars["email"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
		return
	}

	r.ParseForm()

	userToRemove := data.UserDB{
		Email:email,
	}
	err = userToRemove.ReadUUID()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read user UUID.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = projectDB.RemoveEmployee(userToRemove)
	if err != nil {
		utils.Error(err, "Error occurred while trying to remove employee.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
}

func TasksGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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

	tasks, err := projectDB.ReadTasks()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ViewData struct {
		Username  string
		Specialty string
		Employees []data.Employee
		Project   data.ProjectDB
		Tasks     []data.Task
	}

	data := ViewData{
		Username:  userDB.Email,
		Specialty: specialty,
		Project:   projectDB,
		Employees: employees,
		Tasks:     tasks,
	}

	utils.GenerateHTML(w, data, "project_layout", "navbar", "project_tasks")
}

func TasksPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
		return
	}

	r.ParseForm()

	pattern := `.*\S.*`
	taskDescription := r.FormValue("task")

	if matched, err := regexp.Match(pattern, []byte(taskDescription)); !matched || err != nil {
		session.AddFlash("Task description can't be empty!")

		err = session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid+"/tasks", http.StatusSeeOther)
		return
	}

	task := data.Task {
		Description: taskDescription,
		Email:       r.FormValue("employees"),
	}

	err = projectDB.AddTask(task)
	if err != nil {
		utils.Error(err, "Error occurred while trying add task to db.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/project/"+projectDB.Uuid+"/tasks", http.StatusSeeOther)
}

func JoinPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If the user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	if err := userDB.ReadUUID(); err != nil {
		utils.Error(err, "Error occurred while trying to read user UUID.\n")
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
				utils.Error(err, "Error occurred while trying to save session.\n")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return

		} else {
			utils.Error(err, "Error occurred while trying to join project.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/project/"+projectUUID, http.StatusSeeOther)
}

func UpdateProjectDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := usession.Store.Get(r, "session")
	userSession := getUserSession(session)
	userDB := data.UserDB{Email: userSession.Email}

	// If user is unauthenticated, add flash message and return
	if auth := userSession.Authenticated; !auth {
		session.AddFlash("Sign in your account!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	projectUUID := vars["key"]

	// If user tries to get an access to project in which he's not in, he gets an error
	if exists := userDB.IsInProject(projectUUID); !exists {
		session.AddFlash("You have no access to this project!")
		err := session.Save(r, w)

		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	projectDB := data.ProjectDB{Uuid: projectUUID}
	err := projectDB.ReadName()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read project name.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees, err := projectDB.ReadEmployees()
	if err != nil {
		utils.Error(err, "Error occurred while trying to read employees of project.\n")
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
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
		return
	}

	r.ParseForm()

	pattern := `.*\S.*`
	projectDescription := r.FormValue("projectDescription")

	if matched, err := regexp.Match(pattern, []byte(projectDescription)); !matched || err != nil {
		session.AddFlash("Project description can't be empty!")

		err = session.Save(r, w)
		if err != nil {
			utils.Error(err, "Error occurred while trying to save session.\n")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/project/"+projectDB.Uuid+"/tasks", http.StatusSeeOther)
		return
	}

	err = projectDB.UpdateProjectDescription(projectDescription)
	if err != nil {
		utils.Error(err, "Error occurred while trying to update project's description.\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/project/"+projectDB.Uuid, http.StatusSeeOther)
}

func AddReportHandler(w http.ResponseWriter, r *http.Request) {

}