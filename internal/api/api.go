package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/swaggo/swag/cmd/swag@latest init -g api.go -o ../../docs --parseDependency

type API struct {
	service Service
}

func New(s Service) *API {
	return &API{
		service: s,
	}
}

func (a *API) AddRoutes(s *http.ServeMux) {
	s.HandleFunc("/health", a.health)
	s.HandleFunc("POST /users", a.CreateUser)

	s.HandleFunc("GET /users/{id}/tasks", a.ListTasks)
	s.HandleFunc("POST /users/{id}/tasks/start", a.StartTask)
	s.HandleFunc("POST /users/{id}/tasks/{taskID}/end", a.EndTask)
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func (a *API) writeResp(w http.ResponseWriter, r *http.Request, data any) {
	resp := Response{
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		a.internalServerError(w, r, err)
	}
}

func (a *API) writeErr(w http.ResponseWriter, r *http.Request, status int, err error) {
	slog.Error("Request failed", "status", status, "error", err.Error(), "url", r.URL.Path)

	resp := Response{
		Error: err.Error(),
	}

	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to write response", "error", err)
		return
	}
}

func (a *API) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	a.writeErr(w, r, http.StatusBadRequest, err)
}

func (a *API) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	a.writeErr(w, r, http.StatusInternalServerError, err)
}
