package api

import (
	"encoding/json"
	"net/http"
)

type CreateUserRequest struct {
	PassportNumber string `json:"passportNumber"`
}

// CreateUser creates a new user with the given passport number.
// @Summary Create a new user
// @Description Create a new user with the given passport number
// @Tags users
// @Param passportNumber body CreateUserRequest true "Body"
// @Success 201 "User created"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /users [post]
func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.badRequest(w, r, err)
		return
	}

	if err := a.service.CreateUser(r.Context(), req.PassportNumber); err != nil {
		a.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
