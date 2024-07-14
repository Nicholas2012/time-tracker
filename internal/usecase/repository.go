package usecase

import (
	"context"

	"github.com/Nicholas2012/time-tracker/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int) (*models.User, error)

	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, userID, id int) (*models.Task, error)
	ListTasks(ctx context.Context, userID int) ([]models.Task, error)
}
