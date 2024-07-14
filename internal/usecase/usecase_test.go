package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Nicholas2012/time-tracker/internal/models"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_OK(t *testing.T) {
	s, _ := setup(t)
	err := s.CreateUser(context.TODO(), "1234 567890")
	require.NoError(t, err)
}

func TestCreateUser_InvalidPassportNumber(t *testing.T) {
	s, _ := setup(t)
	err := s.CreateUser(context.TODO(), "1234567890")
	require.EqualError(t, err, "invalid passport number, must have at least 2 parts")
}

func TestCreateUser_BadSeries(t *testing.T) {
	s, _ := setup(t)
	err := s.CreateUser(context.TODO(), "abc 567890")
	require.EqualError(t, err, "invalid passport series, must be a number, got: abc")
}

func TestCreateUser_BadNumber(t *testing.T) {
	s, _ := setup(t)
	err := s.CreateUser(context.TODO(), "1234 abc")
	require.EqualError(t, err, "invalid passport number, must be a number, got: abc")
}

func TestStartTask_OK(t *testing.T) {
	s, repo := setup(t)

	testUser := &models.User{ID: 99}

	repo.GetUserFn = func(ctx context.Context, id int) (*models.User, error) {
		return testUser, nil
	}
	repo.CreateTaskFn = func(ctx context.Context, task *models.Task) error {
		task.ID = 1
		require.Equal(t, testUser.ID, task.UserID)
		require.NotEmpty(t, task.Since)
		return nil
	}

	id, err := s.StartTask(context.TODO(), testUser.ID)
	require.NoError(t, err)
	require.NotZero(t, id)
}

func TestStartTask_UserNotFound(t *testing.T) {
	s, repo := setup(t)
	repo.GetUserFn = func(ctx context.Context, id int) (*models.User, error) {
		return nil, sql.ErrNoRows
	}

	_, err := s.StartTask(context.TODO(), 1)

	require.ErrorIs(t, err, ErrNotFound)
}

func TestStartTask_Error(t *testing.T) {
	s, repo := setup(t)

	testUser := &models.User{ID: 99}

	repo.GetUserFn = func(ctx context.Context, id int) (*models.User, error) {
		return testUser, nil
	}
	repo.CreateTaskFn = func(ctx context.Context, task *models.Task) error {
		return errors.New("test error")
	}

	_, err := s.StartTask(context.TODO(), 1)

	require.EqualError(t, err, "create task: test error")
}

func TestEndTask_OK(t *testing.T) {
	s, repo := setup(t)

	now := time.Now().Add(-time.Hour)
	testUser := &models.User{ID: 99}
	testTask := &models.Task{ID: 1, UserID: testUser.ID, Since: now}

	repo.GetUserFn = func(ctx context.Context, id int) (*models.User, error) {
		return testUser, nil
	}
	repo.GetTaskFn = func(ctx context.Context, userID, id int) (*models.Task, error) {
		return testTask, nil
	}
	repo.UpdateTaskFn = func(ctx context.Context, task *models.Task) error {
		require.NotZero(t, task.ID)
		require.NotEmpty(t, task.Since)
		require.NotEmpty(t, task.Until)
		require.NotZero(t, task.Minutes)
		require.Equal(t, 60*time.Minute, task.Until.Sub(task.Since).Truncate(time.Minute))

		return nil
	}

	err := s.EndTask(context.TODO(), testUser.ID, testTask.ID)
	require.NoError(t, err)
}

func TestEndTask_UserNotFound(t *testing.T) {
	s, repo := setup(t)
	repo.GetUserFn = func(ctx context.Context, id int) (*models.User, error) {
		return nil, sql.ErrNoRows
	}

	err := s.EndTask(context.TODO(), 1, 1)

	require.ErrorIs(t, err, ErrNotFound)
}

func TestEndTask_TaskNotFound(t *testing.T) {
	s, repo := setup(t)
	repo.GetUserFn = func(ctx context.Context, id int) (*models.User, error) {
		return &models.User{ID: 99}, nil
	}
	repo.GetTaskFn = func(ctx context.Context, userID, id int) (*models.Task, error) {
		return nil, sql.ErrNoRows
	}

	err := s.EndTask(context.TODO(), 99, 1)

	require.ErrorIs(t, err, ErrNotFound)
}

func TestListTasks_OK(t *testing.T) {
	s, repo := setup(t)

	testUser := &models.User{ID: 99}
	testTasks := []models.Task{
		{ID: 1, UserID: testUser.ID, Since: time.Now(), Until: time.Now().Add(time.Hour), Minutes: 60},
	}

	repo.ListTasksFn = func(ctx context.Context, userID int) ([]models.Task, error) {
		return testTasks, nil
	}

	tasks, err := s.ListTasks(context.TODO(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, testTasks, tasks)
}

func setup(_ *testing.T) (*Service, *repositoryMock) {
	repo := &repositoryMock{}
	return New(repo), repo
}
