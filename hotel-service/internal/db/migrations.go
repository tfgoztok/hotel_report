package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations executes the database migrations from the specified path.
func RunMigrations(db *sql.DB, migrationPath string) error {
	// Create a new Postgres driver instance using the provided database connection.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %v", err)
	}

	// Create a new migration instance with the specified migration path and database driver.
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create the migration instance: %v", err)
	}

	// Run the migrations. If there are no changes, it will not return an error.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %v", err)
	}

	// Return nil if migrations were successful or if there were no changes.
	return nil
}
