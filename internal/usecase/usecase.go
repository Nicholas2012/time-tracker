package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nicholas2012/time-tracker/internal/models"
)

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, passportNumber string) error {
	parts := strings.Split(passportNumber, " ")
	if len(parts) < 2 {
		return errors.New("invalid passport number, must have at least 2 parts")
	}

	series, err := strconv.Atoi(parts[0])
	if err != nil {
		return errors.New("invalid passport series, must be a number, got: " + parts[0])
	}

	number, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("invalid passport number, must be a number, got: " + parts[1])
	}

	newUser := &models.User{
		PassportSerie:  series,
		PassportNumber: number,
	}

	// todo make call to name service to get user name

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Service) StartTask(ctx context.Context, userID int) (int, error) {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("get user: %w", err)
	}

	task := models.NewTask(user.ID)

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}

	return task.ID, nil
}

func (s *Service) EndTask(ctx context.Context, userID, taskID int) error {
	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("get user: %w", err)
	}

	task, err := s.repo.GetTask(ctx, user.ID, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("get task: %w", err)
	}

	task.Until = time.Now()
	task.Minutes = int(task.Until.Sub(task.Since).Minutes())

	if err := s.repo.UpdateTask(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

func (s *Service) ListTasks(ctx context.Context, userID int) ([]models.Task, error) {
	tasks, err := s.repo.ListTasks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return tasks, nil
}
