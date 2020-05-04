package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
	"github.com/lynshi/cuisine-calendar-api/internal/models"
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
		log.Fatal().Err(err).Msg("could not connect to Docker")
	}

	dbname := "testdatabase"
	resource, err := pool.Run("postgres", "9.6", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbname})
	if err != nil {
		log.Fatal().Err(err).Msg("could not start resource")
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), dbname))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}

	defer func() {
		err = pool.Purge(resource)
		if err != nil {
			log.Error().Err(err).Msg("could not purge resource")
		}
	}()

	var gdb *gorm.DB
	gdb, err = gorm.Open("postgres", db)
	if err != nil {
		log.Fatal().Err(err).Msg("could not open connection from Gorm")
	}

	testApp.db = &database.DB{
		DB: gdb,
	}

	testApp.db.AutoMigrate(&models.Recipe{})

	code = m.Run()
}

func TestGetRecipe(t *testing.T) {
	id := 103
	name := "test recipe item"
	servings := 2
	ingredients := json.RawMessage(`{"salt": "1 tbsp"}`)
	created := time.Now().Round(time.Microsecond)
	updated := time.Now().Round(time.Microsecond)
	owner := "me"
	permissions := "everyone"

	dbItem := models.Recipe{
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

	req, err := http.NewRequest("GET", fmt.Sprintf("/getRecipe?id=%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(testApp.router, req)
	checkResponseCode(t, http.StatusOK, response)

	var expectedIngredients map[string]string
	expectedIngredients, err = parseIngredientsJSONB(&dbItem.Ingredients)

	if err != nil {
		t.Fatal(err)
	}

	result := models.GetRecipeResponse{}
	expected := models.GetRecipeResponse{
		RecipeID:    id,
		Name:        name,
		Ingredients: expectedIngredients,
		Servings:    servings,
		CreatedAt:   created,
		UpdatedAt:   updated,
		Owner:       owner,
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(expected, result) {
		t.Errorf(
			"handler returned unexpected body: want %+v got %+v",
			expected, result,
		)
	}
}

func TestGetRecipeNonexistentID(t *testing.T) {
	req, err := http.NewRequest("GET", "/getRecipe?id=42", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(testApp.router, req)
	checkResponseCode(t, http.StatusInternalServerError, response)
}

func TestGetRecipeStringId(t *testing.T) {
	req, err := http.NewRequest("GET", "/getRecipe?id=fail", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(testApp.router, req)
	checkResponseCode(t, http.StatusBadRequest, response)
}

func TestPutRecipeWithoutID(t *testing.T) {
	name := "test recipe item"
	servings := 2
	ingredients := map[string]string{
		"salt": "10 tbsp",
	}
	created := time.Now().Round(time.Microsecond)
	updated := time.Now().Round(time.Microsecond)

	putRecipeRequest := models.PutRecipeRequest{
		Name:        name,
		Servings:    servings,
		Ingredients: ingredients,
		CreatedAt:   created,
		UpdatedAt:   updated,
	}

	jsonStr, err := json.Marshal(putRecipeRequest)
	if err != nil {
		t.Fatal(err)
	}

	var req *http.Request
	req, err = http.NewRequest("POST", "/putRecipe", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(testApp.router, req)
	checkResponseCode(t, http.StatusOK, response)

	var result models.PutRecipeResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	var recipe models.Recipe
	recipe, err = testApp.db.GetRecipeByID(result.RecipeID)
	if err != nil {
		t.Fatal(err)
	}

	if name != recipe.Name {
		t.Errorf("expected recipe name %s got %s", name, recipe.Name)
	}

	if servings != recipe.Servings {
		t.Errorf("expected recipe servings %d got %d", servings, recipe.Servings)
	}

	var parsedIngredients map[string]string
	parsedIngredients, err = parseIngredientsJSONB(&recipe.Ingredients)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(ingredients, parsedIngredients) {
		t.Errorf("expected recipe ingredients %v got %v", ingredients, parsedIngredients)
	}

	if !cmp.Equal(created, recipe.CreatedAt) {
		t.Errorf("expected recipe created at %v got %v", created, recipe.CreatedAt)
	}

	if !cmp.Equal(updated, recipe.UpdatedAt) {
		t.Errorf("expected recipe updated at %v got %v", updated, recipe.UpdatedAt)
	}
}
