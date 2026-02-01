package common_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodHead, "/healthcheck", nil)
	w := httptest.NewRecorder()
	common.HealthCheck(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	assert.Equal(t, "204 No Content", res.Status)
}

func TestPing(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	common.Ping(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	assert.Equal(t, "200 OK", res.Status)

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data struct {
		Message string `json:"message"`
	}

	err = json.Unmarshal(bodyBytes, &data)
	require.NoError(t, err)

	assert.Equal(t, "pong", data.Message)
}
