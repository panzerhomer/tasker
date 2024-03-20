package handlers

import (
	"encoding/json"
	"net/http"
	"rest/internal/domain"
	"rest/internal/services"
)

type TaskHandler struct {
	taskServo services.TaskService
}

func NewTaskHandler(taskServo services.TaskService) *TaskHandler {
	return &TaskHandler{taskServo: taskServo}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		response(w, failed, http.StatusInternalServerError)
		return
	}

	// projectNameParam := chi.URLParam(r, "name")

	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		response(w, "invalid data", http.StatusBadRequest)
		return
	}

	if err := task.Validate(); err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.taskServo.CreateTask(ctx, userID, &task)
	if err != nil {
		response(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response(w, success, http.StatusCreated)
}
