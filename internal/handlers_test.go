package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

var testApp App

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	testApp.initializeRouter()
	os.Exit(m.Run())
}

func createMockDB(t *testing.T) *sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	var gdb *gorm.DB
	gdb, err = gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock GORM connection", err)
	}

	testApp.db = gdb

	t.Cleanup(func() {
		gdb.Close()
		db.Close()
		testApp.db = nil
	})

	return &mock
}

func TestGetRecipe(t *testing.T) {
	createMockDB(t) // mock = createMockDB(t)

	req, err := http.NewRequest("GET", "/recipe/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

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
