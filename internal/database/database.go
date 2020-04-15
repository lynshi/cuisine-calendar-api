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
func InitializeDatabaseConnection(dbname, user, password, host string, port int) *DB {
	log.Info().Msg(fmt.Sprintf("Connecting to database %v", dbname))

	connectionString :=
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d", dbname, user, password, host, port)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	return &DB{db}
}

func (db *DB) GetRecipeByID(id int) *Recipe {
	var recipe Recipe
	db.First(&recipe, id)
	return &recipe
}
