package api

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Nicholas2012/time-tracker/internal/models"
	"github.com/stretchr/testify/require"
)

func TestTasksList_OK(t *testing.T) {
	srv, sm := setup(t)

	sm.listTasksFn = func(_ context.Context, userID int) ([]models.Task, error) {
		return []models.Task{
			{
				ID:      81,
				UserID:  51,
				Since:   time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC),
				Until:   time.Date(2021, 10, 1, 1, 0, 0, 0, time.UTC),
				Minutes: 60,
			},
			{
				ID:      82,
				UserID:  51,
				Since:   time.Date(2021, 10, 1, 2, 0, 0, 0, time.UTC),
				Until:   time.Date(2021, 10, 1, 3, 0, 0, 0, time.UTC),
				Minutes: 60,
			},
		}, nil
	}

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/users/51/tasks", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.JSONEq(t, `{"data": [{"id": 81, "since": "2021-10-01T00:00:00Z", "until": "2021-10-01T01:00:00Z", "minutes": 60}, {"id": 82, "since": "2021-10-01T02:00:00Z", "until": "2021-10-01T03:00:00Z", "minutes": 60}]}`, string(body))
}
