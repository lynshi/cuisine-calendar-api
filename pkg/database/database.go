package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // imported but not used
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DB is a wrapper for a gorm.DB object.
type DB struct {
	*gorm.DB
}

// New opens a connection to the given database.
func New(dbname, user, password, host string, port int, ssl bool) (*DB, error) {
	log.Info().Msg(fmt.Sprintf("Connecting to database %v", dbname))

	connectionString :=
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d", dbname, user, password, host, port)
	if !ssl {
		connectionString = fmt.Sprintf("%s sslmode=disable", connectionString)
	}

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		err = errors.Wrap(err, "could not connect to database")
		return nil, err
	}

	return &DB{db}, nil
}
