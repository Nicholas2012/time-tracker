package api

import (
	"context"

	"github.com/Nicholas2012/time-tracker/internal/models"
)

type Service interface {
	CreateUser(ctx context.Context, passportNumber string) error

	StartTask(ctx context.Context, userID int) (int, error)
	EndTask(ctx context.Context, userID, taskID int) error
	ListTasks(ctx context.Context, userID int) ([]models.Task, error)
}
