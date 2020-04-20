package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"

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

	// Database is up, so we can connect using our function instead.
	db.Close()

	var port int
	port, err = strconv.Atoi(resource.GetPort("5432/tcp"))
	if err != nil {
		log.Fatal().Err(err).Msg("Could not convert port to int")
	}

	testDB = InitializeDatabaseConnection(dbname, "postgres", "secret", "localhost", port, false)
	testDB.createTables()

	code = m.Run()
}

var tableExistsTests = []struct {
	name string
	in   interface{}
	out  bool
}{
	{"Recipe table", &Recipe{}, true},
	{"TextInstruction table", &TextInstruction{}, true},
}

func TestTablesCreated(t *testing.T) {
	for _, tt := range tableExistsTests {
		t.Run(tt.name, func(t *testing.T) {
			if testDB.HasTable(tt.in) != tt.out {
				t.Errorf("Could not find table in database")
			}
		})
	}
}
