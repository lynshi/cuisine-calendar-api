package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testApp App

func TestMain(m *testing.M) {
	testApp.initialize()
	os.Exit(m.Run())
}

func TestGetRecipe(t *testing.T) {
	req, err := http.NewRequest("GET", "/recipe/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// Check the response body is what we expect.
	result := getRecipeResponse{}
	expected := &getRecipeResponse{
		RecipeId: 1,
		Name:     "test",
	}

	err = json.Unmarshal(response.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != *expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expected)
	}
}

func TestGetRecipeStringId(t *testing.T) {
	req, err := http.NewRequest("GET", "/recipe/fail", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testApp.router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
