package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

func executeRequest(r *router.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
