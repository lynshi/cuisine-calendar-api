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

func initializeDatabaseConnection(dbname, user, password, host string, port int) DB {
	log.Info().Msg(fmt.Sprintf("Connecting to database %v", dbname))

	connectionString :=
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d", dbname, user, password, host, port)
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return DB{db}
}
