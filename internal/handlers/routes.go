package handlers

import (
	authmid "rest/internal/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Routes(userHandler *UserHandler, taskHandler *TaskHandler, projectHandler *ProjectHandler) chi.Router {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(middleware.RequestID)
	root.Post("/login", userHandler.Login)
	root.Post("/signin", userHandler.Signin)

	r := chi.NewRouter()
	r.Use(authmid.Auth)
	r.Get("/hello", userHandler.GetHello)
	r.Get("/logout", userHandler.Logout)

	r.Post("/projects", projectHandler.CreateProject)
	r.Get("/projects", projectHandler.GetAllProjects)
	r.Put("/projects/{name}", projectHandler.JoinProject)
	// r.Get("/projects/{name}/tasks", projectHandler.GetAllProjectTasks)

	root.Mount("/api", r)

	return root
}
