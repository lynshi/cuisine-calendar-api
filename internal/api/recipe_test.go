package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lynshi/cuisine-calendar-api/internal/database"
	"github.com/lynshi/cuisine-calendar-api/internal/router"
)

var (
	testApp appContext
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	testApp.debug = true
	testApp.router = router.NewRouter()
	testApp.setupRouter()

	code := 0
	defer func() {
		os.Exit(code)
	}()

	var db *sql.DB
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to Docker")
	}

	dbname := "testdatabase"
	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbname})
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start resource")
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), dbname))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.Error().Err(err).Msg("Could not purge resource")
		}
	}()

	var gdb *gorm.DB
	gdb, err = gorm.Open("postgres", db)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not open connection from Gorm")
	}

	testApp.db = &database.DB{
		DB: gdb,
	}

	testApp.db.AutoMigrate(&database.Recipe{})

	code = m.Run()
}

func TestGetRecipe(t *testing.T) {
	id := 1
	name := "test recipe item"
	servings := 2
	ingredients := json.RawMessage(`{"salt": 1}`)
	created := time.Now()
	updated := time.Now()
	owner := "me"
	permissions := "everyone"

	dbItem := database.Recipe{
		ID:          id,
		Name:        name,
		Servings:    servings,
		Ingredients: postgres.Jsonb{RawMessage: ingredients},
		CreatedAt:   created,
		UpdatedAt:   updated,
		Owner:       owner,
		Permissions: permissions,
	}
	testApp.db.Create(&dbItem)

	req, err := http.NewRequest("GET", "/recipe/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var expectedIngredients map[string]int
	expectedIngredients, err = parseIngredientsJSONB(
		&postgres.Jsonb{RawMessage: ingredients},
	)

	if err != nil {
		t.Fatal(err)
	}

	result := getRecipeResponse{}
	expected := getRecipeResponse{
		RecipeID:    id,
		Name:        name,
		Ingredients: expectedIngredients,
		Servings:    servings,
		CreatedAt:   created,
		UpdatedAt:   updated,
		Owner:       owner,
	}

	err = json.Unmarshal(response.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	if cmp.Equal(result, expected) {
		t.Errorf(
			"handler returned unexpected body: got %v want %v",
			result, expected,
		)
	}
}

func TestGetRecipeStringId(t *testing.T) {
	req, err := http.NewRequest("GET", "/recipe/fail", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
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
