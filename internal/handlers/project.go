package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rest/internal/domain"

	"github.com/go-chi/chi"
)

type ProjectService interface {
	CreateProject(ctx context.Context, authorID int64, project domain.Project) error
	JoinProject(ctx context.Context, projectName string, userID int64, role int8) error
	GetAllProjects(ctx context.Context, userID int64) ([]*domain.Project, error)
	LeaveProject(ctx context.Context, projectName string, userID int64) error
	CreateProjectTask(ctx context.Context, userID int64, projectName string, task domain.Task) error
	GetProjectUsers(ctx context.Context, userID int64, projectName string) ([]*domain.ProjectUser, error)
}

type ProjectHandler struct {
	projectServo ProjectService
}

func NewProjectHandler(projectServo ProjectService) *ProjectHandler {
	return &ProjectHandler{projectServo: projectServo}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		response(w, "invalid data", http.StatusBadRequest)
		return
	}

	if err := project.Validate(); err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.projectServo.CreateProject(ctx, userID, project)
	if err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}

	response(w, success, http.StatusCreated)
}

func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projects, err := h.projectServo.GetAllProjects(ctx, userID)
	if err != nil {
		response(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, projects, http.StatusOK)
}

// here
func (h *ProjectHandler) GetAllProjectUsers(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projectNameParam := chi.URLParam(r, "name")
	users, err := h.projectServo.GetProjectUsers(ctx, userID, projectNameParam)
	if err != nil {
		response(w, "unable to get project users", http.StatusInternalServerError)
		return
	}

	if users != nil {
		responseJSON(w, users, http.StatusOK)
		return
	}

	response(w, failed, http.StatusInternalServerError)

}

func (h *ProjectHandler) GetAllProjectTasks(w http.ResponseWriter, r *http.Request) {
	dateParam := chi.URLParam(r, "name")
	log.Println(dateParam)
}

func (h *ProjectHandler) JoinProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projectNameParam := chi.URLParam(r, "name")
	err := h.projectServo.JoinProject(ctx, projectNameParam, userID, domain.UserRole)
	if err != nil {
		response(w, failed, http.StatusInternalServerError)
		return
	}

	response(w, success, http.StatusOK)
}

func (h *ProjectHandler) LeaveProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projectNameParam := chi.URLParam(r, "name")

	err := h.projectServo.LeaveProject(ctx, projectNameParam, userID)
	if err != nil {
		response(w, failed, http.StatusInternalServerError)
		return
	}

	response(w, success, http.StatusOK)
}

func (h *ProjectHandler) CreateProjectTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response(w, "invalid data", http.StatusBadRequest)
		return
	}

	if err := task.Validate(); err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectNameParam := chi.URLParam(r, "name")
	err = h.projectServo.CreateProjectTask(ctx, userID, projectNameParam, task)
	if err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}
	response(w, success, http.StatusOK)
}
