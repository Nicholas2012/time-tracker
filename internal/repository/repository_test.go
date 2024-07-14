package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"os"
	"testing"
	"time"

	"github.com/Nicholas2012/time-tracker/internal/models"
	"github.com/Nicholas2012/time-tracker/pkg/database"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("DOCKER_HOST", "unix:///Users/ok/.colima/default/docker.sock")
}

func TestRepository(t *testing.T) {
	repo := setup(t)

	user := &models.User{
		Name:           "John",
		Surname:        "Doe",
		Patronymic:     "Smith",
		PassportSerie:  1234,
		PassportNumber: 567890,
	}

	t.Run("CreateUser", func(t *testing.T) {
		err := repo.CreateUser(context.Background(), user)
		require.NoError(t, err)
		require.NotZero(t, user.ID)

		t.Run("GetUser", func(t *testing.T) {
			u, err := repo.GetUser(context.Background(), user.ID)
			require.NoError(t, err)
			require.Equal(t, user.ID, u.ID)
			require.Equal(t, user.Name, u.Name)
			require.Equal(t, user.Surname, u.Surname)
			require.Equal(t, user.Patronymic, u.Patronymic)
			require.Equal(t, user.PassportSerie, u.PassportSerie)
			require.Equal(t, user.PassportNumber, u.PassportNumber)
		})
	})

	t.Run("CreateTask", func(t *testing.T) {
		task := &models.Task{
			UserID:  user.ID,
			Since:   time.Now().Add(-time.Hour),
			Until:   time.Now(),
			Minutes: 60,
		}
		err := repo.CreateTask(context.Background(), task)
		require.NoError(t, err)
		require.NotZero(t, task.ID)

		t.Run("UpdateTask", func(t *testing.T) {
			task.Minutes = 120
			err := repo.UpdateTask(context.Background(), task)
			require.NoError(t, err)
		})

		t.Run("GetTask", func(t *testing.T) {
			result, err := repo.GetTask(context.Background(), user.ID, task.ID)
			require.NoError(t, err)
			require.Equal(t, task.ID, result.ID)
			require.Equal(t, task.UserID, result.UserID)
			require.Equal(t, task.Since.Unix(), result.Since.Unix())
			require.Equal(t, task.Until.Unix(), result.Until.Unix())
			require.Equal(t, 120, result.Minutes)
		})

		t.Run("ListTasks", func(t *testing.T) {
			tasks, err := repo.ListTasks(context.Background(), user.ID)
			require.NoError(t, err)
			require.NotEmpty(t, tasks)

			require.Len(t, tasks, 1)
			require.Equal(t, task.ID, tasks[0].ID)
			require.Equal(t, task.UserID, tasks[0].UserID)
			require.Equal(t, task.Since.Unix(), tasks[0].Since.Unix())
			require.Equal(t, task.Until.Unix(), tasks[0].Until.Unix())
			require.Equal(t, 120, tasks[0].Minutes)
		})
	})
}

func TestMigrations(t *testing.T) {
	setup(t)
}

func setup(t *testing.T) *Repository {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	// Ping docker
	if err := pool.Client.Ping(); err != nil {
		t.Skipf("test will be skipped because docker is not running, ping error: %s", err)
	}

	// Run postgres
	resource, err := pool.Run("postgres", "latest",
		[]string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_DB=testdb",
		})
	require.NoError(t, err)

	var db *sql.DB

	// Wait for postgres to be ready
	err = pool.Retry(func() error {
		newDB, err := database.New(fmt.Sprintf("postgres://postgres:secret@localhost:%s/testdb?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		db = newDB
		return nil
	})
	require.NoError(t, err)

	// Set purge function
	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource))
	})

	// Apply migrations
	require.NoError(t, database.ApplyMigrations(db))

	return New(db)
}
