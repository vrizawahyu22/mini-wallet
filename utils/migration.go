package utils

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrationPool(db *sql.DB, config *BaseConfig, isSkipForTesting ...bool) {
	defer func() {
		err := db.Close()
		LogAndPanicIfError(err, "failed when closing database (migration)")
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	LogAndPanicIfError(err, "cannot create new migrate instance")

	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationURL,
		"postgres", driver)
	LogAndPanicIfError(err, "failed to create postgres instance")

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		LogAndPanicIfError(err, "failed to run migrate up")
	}

	LogInfo("db migrated successfully")
}
