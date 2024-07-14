package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTasksEnd_OK(t *testing.T) {
	srv, sm := setup(t)

	sm.endTaskFn = func(ctx context.Context, userID, taskID int) error {
		require.Equal(t, 51, userID)
		require.Equal(t, 69, taskID)
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/users/51/tasks/69/end", nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
}
