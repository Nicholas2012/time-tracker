package api

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTasksStart_OK(t *testing.T) {
	s, sm := setup(t)

	sm.startTaskFn = func(ctx context.Context, userID int) (int, error) {
		require.Equal(t, 51, userID)
		return 69, nil
	}

	req, err := http.NewRequest(http.MethodPost, s.URL+"/users/51/tasks/start", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.JSONEq(t, `{"data": {"task_id": 69}}`, string(body))
}
