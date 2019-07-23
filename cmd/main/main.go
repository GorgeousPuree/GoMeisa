package main

import (
	_ "Gomeisa/init"
	"Gomeisa/internal/api"
	"github.com/gorilla/mux"
	"os"

	"log"
	"net/http"
)

func main() {
	port := os.Getenv("SERVICE_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	r := mux.NewRouter()
	projectRouter := r.PathPrefix("/project/{key}").Subrouter()
	inviteRouter := projectRouter.PathPrefix("/invite").Subrouter()
	taskRouter := projectRouter.PathPrefix("/tasks").Subrouter()
	//removeRouter := projectRouter.PathPrefix("/removeEmployee/{email}").Subrouter()

	r.HandleFunc("/signin", api.SigninGetHandler).Methods("GET")
	r.HandleFunc("/signin", api.SigninPostHandler).Methods("POST")
	r.HandleFunc("/logout", api.LogoutPostHandler).Methods("POST")

	r.HandleFunc("/signup", api.SignupGetHandler).Methods("GET")
	r.HandleFunc("/signup", api.SignupPostHandler).Methods("POST")
	r.HandleFunc("/join/{key:[^ ]+}/",api.JoinPostHandler)

	r.HandleFunc("/projects", api.ProjectsGetHandler).Methods("GET")
	r.HandleFunc("/createProject", api.ProjectPostHandler).Methods("POST")

	projectRouter.HandleFunc("", api.ProjectGetHandler).Methods("GET")
	//removeRouter.HandleFunc("", removeEmployeeHandler).Methods("POST")
	// Those functions are POST-methods.
	projectRouter.HandleFunc("/removeEmployee/{email}", api.RemoveEmployeeHandler)
	// UPDATE method
	projectRouter.HandleFunc("/updateProjectDescription", api.UpdateProjectDescriptionHandler).Methods("POST")
	projectRouter.HandleFunc("/addReportDescription", api.AddReportHandler).Methods("POST")

	inviteRouter.HandleFunc("", api.InviteGetHandler).Methods("GET")
	inviteRouter.HandleFunc("", api.InvitePostHandler).Methods("POST")

	taskRouter.HandleFunc("", api.TasksGetHandler).Methods("GET")
	taskRouter.HandleFunc("", api.TasksPostHandler).Methods("POST")

	r.HandleFunc("/", api.ProjectsGetHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("../web/static/"))))

	log.Fatal(http.ListenAndServe(":"+port, r))
}