package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lynshi/cuisine-calendar-api/internal/apimodels"
	"github.com/lynshi/cuisine-calendar-api/internal/dbmodels"
	"github.com/lynshi/cuisine-calendar-api/internal/logsetup"	
	"github.com/lynshi/cuisine-calendar-api/pkg/database"
	"github.com/lynshi/cuisine-calendar-api/pkg/router"
)

var (
	testApp *App
)

func TestMain(m *testing.M) {
	logsetup.SetupZerolog(true)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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

	// Database is up, so we can connect using our function instead.
	db.Close()

	var port int
	port, err = strconv.Atoi(resource.GetPort("5432/tcp"))
	if err != nil {
		log.Fatal().Err(err).Msg("could not convert port to int")
	}

	var testDB *database.DB
	testDB, err = database.New(dbname, "postgres", "secret", "localhost", port, false)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create database")
	}
	defer testDB.Close()

	dbmodels.CreateRecipeTable(testDB)

	router := router.New()

	testApp = New(testDB, router)
	testApp.SetupRouter()

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

	dbItem := dbmodels.Recipe{
		ID:          id,
		Name:        name,
		Servings:    servings,
		Ingredients: postgres.Jsonb{RawMessage: ingredients},
		CreatedAt:   created,
		UpdatedAt:   updated,
		Owner:       owner,
		Permissions: permissions,
	}
	dbmodels.AddRecipe(testApp.db, &dbItem)

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

	result := apimodels.GetRecipeResponse{}
	expected := apimodels.GetRecipeResponse{
		RecipeID:    dbItem.ID,
		Name:        dbItem.Name,
		Ingredients: expectedIngredients,
		Servings:    dbItem.Servings,
		CreatedAt:   dbItem.CreatedAt,
		UpdatedAt:   dbItem.UpdatedAt,
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

	putRecipeRequest := apimodels.PutRecipeRequest{
		Name:        name,
		Servings:    servings,
		Ingredients: ingredients,
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

	var result apimodels.PutRecipeResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	var recipe dbmodels.Recipe
	recipe, err = dbmodels.GetRecipeByID(testApp.db, result.RecipeID)
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
}

func TestPutRecipeUpdatesExisting(t *testing.T) {
	name := "test recipe item 34"
	servings := 8
	ingredients := map[string]string{
		"salt": "5 tbsp",
	}

	putRecipeRequest := apimodels.PutRecipeRequest{
		Name:        name,
		Servings:    servings,
		Ingredients: ingredients,
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

	var result apimodels.PutRecipeResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	recipeID := result.RecipeID

	name = "totally new name"
	servings = 21
	ingredients = map[string]string{
		"salt": "not enough",
	}

	putRecipeRequest = apimodels.PutRecipeRequest{
		ID:          &recipeID,
		Name:        name,
		Servings:    servings,
		Ingredients: ingredients,
	}

	jsonStr, err = json.Marshal(putRecipeRequest)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", "/putRecipe", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	response = executeRequest(testApp.router, req)
	checkResponseCode(t, http.StatusOK, response)

	var recipe dbmodels.Recipe
	recipe, err = dbmodels.GetRecipeByID(testApp.db, result.RecipeID)
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
}

func TestPutRecipeMakesNewID(t *testing.T) {
	name := "test recipe item"
	servings := 2
	ingredients := map[string]string{
		"salt": "10 tbsp",
	}

	putRecipeRequest := apimodels.PutRecipeRequest{
		Name:        name,
		Servings:    servings,
		Ingredients: ingredients,
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

	var result apimodels.PutRecipeResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	firstID := result.RecipeID

	req, err = http.NewRequest("POST", "/putRecipe", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	response = executeRequest(testApp.router, req)
	checkResponseCode(t, http.StatusOK, response)

	decoder = json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if firstID == result.RecipeID {
		t.Errorf("recipe ID is the same for both requests")
	}
}
