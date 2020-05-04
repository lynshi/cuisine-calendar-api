package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func checkTimestampOccursAfter(t *testing.T, expected time.Time, actual time.Time) {
	if !expected.Equal(actual) && !expected.Before(actual) {
		t.Errorf("expected timestamp %v occurs after actual %v", expected, actual)
	}
}
