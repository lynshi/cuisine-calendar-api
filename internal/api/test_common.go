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

func checkResponseCode(t *testing.T, expected int, response *httptest.ResponseRecorder) {
	if expected != response.Code {
		t.Errorf("expected response code %d. Got %d from %+v\n", expected, response.Code, response)
	}
}
