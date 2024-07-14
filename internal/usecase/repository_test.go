package usecase

import (
	"context"

	"github.com/Nicholas2012/time-tracker/internal/models"
)

type repositoryMock struct {
	CreateUserFn func(ctx context.Context, user *models.User) error
	GetUserFn    func(ctx context.Context, id int) (*models.User, error)
	CreateTaskFn func(ctx context.Context, task *models.Task) error
	UpdateTaskFn func(ctx context.Context, task *models.Task) error
	GetTaskFn    func(ctx context.Context, userID, id int) (*models.Task, error)
	ListTasksFn  func(ctx context.Context, userID int) ([]models.Task, error)
}

func (r *repositoryMock) CreateUser(ctx context.Context, user *models.User) error {
	if r.CreateUserFn == nil {
		return nil
	}
	return r.CreateUserFn(ctx, user)
}

func (r *repositoryMock) GetUser(ctx context.Context, id int) (*models.User, error) {
	if r.GetUserFn == nil {
		return nil, nil
	}
	return r.GetUserFn(ctx, id)
}

func (r *repositoryMock) CreateTask(ctx context.Context, task *models.Task) error {
	if r.CreateTaskFn == nil {
		return nil
	}
	return r.CreateTaskFn(ctx, task)
}

func (r *repositoryMock) UpdateTask(ctx context.Context, task *models.Task) error {
	if r.UpdateTaskFn == nil {
		return nil
	}
	return r.UpdateTaskFn(ctx, task)
}

func (r *repositoryMock) GetTask(ctx context.Context, userID, id int) (*models.Task, error) {
	if r.GetTaskFn == nil {
		return nil, nil
	}
	return r.GetTaskFn(ctx, userID, id)
}

func (r *repositoryMock) ListTasks(ctx context.Context, userID int) ([]models.Task, error) {
	if r.ListTasksFn == nil {
		return nil, nil
	}
	return r.ListTasksFn(ctx, userID)
}
