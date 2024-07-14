package api

import (
	"errors"
	"net/http"
	"strconv"
)

// EndTask ends a task for a user
// @Summary End a task
// @Tags tasks
// @Param id path number true "User ID"
// @Param taskID path number true "Task ID"
// @Success 200 "Task started"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users/{id}/tasks/{taskID}/end [post]
func (a *API) EndTask(w http.ResponseWriter, r *http.Request) {
	userStr := r.PathValue("id")
	if userStr == "" {
		a.badRequest(w, r, errors.New("missing user ID"))
		return
	}

	userID, err := strconv.Atoi(userStr)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	taskStr := r.PathValue("taskID")
	if taskStr == "" {
		a.badRequest(w, r, errors.New("missing task ID"))
		return
	}

	taskID, err := strconv.Atoi(taskStr)
	if err != nil {
		a.badRequest(w, r, err)
		return
	}

	if err := a.service.EndTask(r.Context(), userID, taskID); err != nil {
		a.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
