package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	testDB *DB
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
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

	testDB = InitializeDatabaseConnection(dbname, "postgres", "secret", "localhost", port, false)
	testDB.createTables()

	code = m.Run()
}

func createRecipeChannel() <-chan models.Recipe {
	recipeChannel := make(chan models.Recipe)

	go func() {
		defer close(recipeChannel)

		recipeID := 5
		recipeName := "database test recipe item"
		recipeServings := 3
		recipeIngredients := json.RawMessage(`{"salt": "1 tbsp"}`)
		recipeOwner := "me"
		recipePermissions := "everyone"
		for {
			recipeChannel <- models.Recipe{
				ID:          recipeID,
				Name:        recipeName,
				Servings:    recipeServings,
				Ingredients: postgres.Jsonb{RawMessage: recipeIngredients},
				Owner:       recipeOwner,
				Permissions: recipePermissions,
			}
			recipeID++
		}
	}()

	return recipeChannel
}

var tableExistsTests = []struct {
	name string
	in   interface{}
	out  bool
}{
	{"Recipe table", &models.Recipe{}, true},
	{"TextInstruction table", &models.TextInstruction{}, true},
}

func TestTablesCreated(t *testing.T) {
	for _, tt := range tableExistsTests {
		t.Run(tt.name, func(t *testing.T) {
			if testDB.HasTable(tt.in) != tt.out {
				t.Errorf("could not find table in database")
			}
		})
	}
}

var (
	// recipeComparer compares Recipes but ignores time.Time fields
	// which are generated by the database.
	recipeComparer = cmp.Comparer(func(x, y models.Recipe) bool {
		if x.ID != y.ID {
			return false
		} else if x.Name != y.Name {
			return false
		} else if x.Servings != y.Servings {
			return false
		} else if !cmp.Equal(x.Ingredients, y.Ingredients) {
			return false
		} else if x.Owner != y.Owner {
			return false
		} else if x.Permissions != y.Permissions {
			return false
		}

		return true
	})

	recipes = createRecipeChannel()
)

func TestAddRecipe(t *testing.T) {
	recipe := <-recipes
	testDB.AddRecipe(&recipe)

	var result models.Recipe
	testDB.Raw("SELECT * FROM recipes WHERE id = ?", recipe.ID).Scan(&result)

	if !cmp.Equal(recipe, result, recipeComparer) {
		t.Errorf(
			"handler returned unexpected body: want %+v got %+v",
			recipe, result,
		)
	}
}

func TestUpdateRecipe(t *testing.T) {
	recipe := <-recipes

	testDB.AddRecipe(&recipe)

	recipe.Name = "this is a totally new name!!!"
	testDB.UpdateRecipe(&recipe)

	var result models.Recipe
	testDB.Raw("SELECT * FROM recipes WHERE id = ?", recipe.ID).Scan(&result)

	if !cmp.Equal(recipe, result, recipeComparer) {
		t.Errorf(
			"handler returned unexpected body: want %+v got %+v",
			recipe, result,
		)
	}
}

func TestGetRecipeByID(t *testing.T) {
	recipe := <-recipes
	recipeID := recipe.ID
	testDB.Exec("INSERT INTO recipes (id) VALUES (?)", recipeID)

	result, err := testDB.GetRecipeByID(recipeID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if recipeID != result.ID {
		t.Errorf("expected recipe ID %d, got %d", recipeID, result.ID)
	}
}

func TestGetRecipeByIDNonexistentID(t *testing.T) {
	_, err := testDB.GetRecipeByID(100000)
	if err == nil {
		t.Errorf("expected \"record not found error\"")
	}
}
