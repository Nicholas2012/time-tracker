package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type ListTasksResponse []Task

type Task struct {
	ID      int       `json:"id"`
	Since   time.Time `json:"since"`
	Until   time.Time `json:"until"`
	Minutes int       `json:"minutes"`
}

// ListTasks lists all tasks for the given user.
// @Summary List all tasks for a user
// @Tags tasks
// @Param id path number true "User ID"
// @Success 200 {object} Response{data=ListTasksResponse} "Task started"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users/{id}/tasks/ [get]
func (a *API) ListTasks(w http.ResponseWriter, r *http.Request) {
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

	tasks, err := a.service.ListTasks(r.Context(), userID)
	if err != nil {
		a.internalServerError(w, r, err)
		return
	}

	tasksItems := make([]Task, len(tasks))
	for i, t := range tasks {
		tasksItems[i] = Task{
			ID:      t.ID,
			Since:   t.Since,
			Until:   t.Until,
			Minutes: t.Minutes,
		}
	}

	a.writeResp(w, r, ListTasksResponse(tasksItems))
}
