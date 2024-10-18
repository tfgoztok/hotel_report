package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	// Attempt to open a connection to the database using the provided URL
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		// Return nil and the error if the connection could not be established
		return nil, err
	}

	// Ping the database to verify that the connection is alive
	if err := db.Ping(); err != nil {
		// Return nil and the error if the ping fails
		return nil, err
	}

	// Return the database connection if successful
	return db, nil
}
