package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nicholas2012/time-tracker/internal/models"
)

type serviceMock struct {
	createUserFn func(ctx context.Context, passportNumber string) error
	startTaskFn  func(ctx context.Context, userID int) (int, error)
	endTaskFn    func(ctx context.Context, userID, taskID int) error
	listTasksFn  func(ctx context.Context, userID int) ([]models.Task, error)
}

func (m *serviceMock) CreateUser(ctx context.Context, passportNumber string) error {
	return m.createUserFn(ctx, passportNumber)
}

func (m *serviceMock) StartTask(ctx context.Context, userID int) (int, error) {
	return m.startTaskFn(ctx, userID)
}

func (m *serviceMock) EndTask(ctx context.Context, userID, taskID int) error {
	return m.endTaskFn(ctx, userID, taskID)
}

func (m *serviceMock) ListTasks(ctx context.Context, userID int) ([]models.Task, error) {
	return m.listTasksFn(ctx, userID)
}

func setup(t *testing.T) (*httptest.Server, *serviceMock) {
	mux := http.NewServeMux()
	sm := &serviceMock{}

	New(sm).AddRoutes(mux)

	srv := httptest.NewServer(mux)

	return srv, sm
}
