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

	log.Println("[handler project]", userID, "|", project, "|")

	err := h.projectServo.CreateProject(ctx, userID, project)
	// log.Println("[handler project]", userID, "|", project, "|", err.Error())
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

func (h *ProjectHandler) GetAllProjectUsers(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// dateParam := chi.URLParam(r, "date")
}

func (h *ProjectHandler) GetAllProjectTasks(w http.ResponseWriter, r *http.Request) {
	dateParam := chi.URLParam(r, "date")
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
