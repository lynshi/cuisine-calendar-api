package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog/log"
)

type DB struct {
	*gorm.DB
}

// InitializeDatabaseConnection opens a connection to the given database.
func InitializeDatabaseConnection(dbname, user, password, host string, port int, ssl bool) *DB {
	log.Info().Msg(fmt.Sprintf("Connecting to database %v", dbname))

	connectionString :=
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d", dbname, user, password, host, port)
	if !ssl {
		connectionString = fmt.Sprintf("%s sslmode=disable", connectionString)
	}

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	return &DB{db}
}

func BuildDatabase(dbname, user, password, host string, port int, debug bool) {
	db := InitializeDatabaseConnection(dbname, user, password, host, port, false)
	db.createTables()
}

func (db *DB) createTables() {
	db.DB.AutoMigrate(&Recipe{})
	db.DB.AutoMigrate(&TextInstruction{})
}

func (db *DB) GetRecipeByID(id int) *Recipe {
	var recipe Recipe
	db.First(&recipe, id)
	return &recipe
}
