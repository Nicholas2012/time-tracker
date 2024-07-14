package api

import (
	"errors"
	"net/http"
	"strconv"
)

type StartTaskResponse struct {
	TaskID int `json:"task_id"`
}

// StartTask creates a new task for the given user
// @Summary Create a new task and start it
// @Tags tasks
// @Param id path number true "User ID"
// @Success 200 {object} Response{data=api.StartTaskResponse} "Task started"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users/{id}/tasks [post]
func (a *API) StartTask(w http.ResponseWriter, r *http.Request) {
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

	id, err := a.service.StartTask(r.Context(), userID)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

	a.writeResp(w, r, StartTaskResponse{TaskID: id})
}
