package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // imported but not used
	"github.com/rs/zerolog/log"

	"github.com/lynshi/cuisine-calendar-api/internal/models"
)

// DB is a wrapper for a gorm.DB object.
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

// BuildDatabase creates tables for each modeled entity.
func BuildDatabase(dbname, user, password, host string, port int) {
	log.Info().Msg("Building database")

	db := InitializeDatabaseConnection(dbname, user, password, host, port, false)
	defer db.Close()

	db.createTables()
}

func (db *DB) createTables() {
	db.DB.AutoMigrate(&models.Recipe{})
	db.DB.AutoMigrate(&models.TextInstruction{})
}

// AddRecipe adds an entry in the recipes table using `recipe`.
func (db *DB) AddRecipe(recipe *models.Recipe) {
	db.Create(recipe)
}

// UpdateRecipe updates an existing recipe.
func (db *DB) UpdateRecipe(recipe *models.Recipe) {
	db.Save(recipe)
}

// GetRecipeByID retrieves an entry in recipes by ID, and returns an error if no such entry exists.
func (db *DB) GetRecipeByID(id int) (models.Recipe, error) {
	var recipe models.Recipe
	err := db.First(&recipe, id).Error
	return recipe, err
}
