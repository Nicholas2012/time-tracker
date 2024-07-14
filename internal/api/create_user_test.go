package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser_OK(t *testing.T) {
	srv, sm := setup(t)

	body := `{ "passportNumber": "1234567890" }`

	sm.createUserFn = func(_ context.Context, _ string) error {
		return nil
	}

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/users", strings.NewReader(body))
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreateUser_BadRequest(t *testing.T) {
	srv, _ := setup(t)

	body := `{ "passportNumber": 1234567890 }`

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/users", strings.NewReader(body))
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	var resp Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	require.NoError(t, err)

	require.NotEmpty(t, resp.Error)
	require.Contains(t, resp.Error, "json: cannot unmarshal number into Go struct field CreateUserRequest.passportNumber of type string")
}
