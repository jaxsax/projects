package botv2

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ConnectDB establishes a connection to a postgres database
func ConnectDB(conf *DBConfig) (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Name,
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
