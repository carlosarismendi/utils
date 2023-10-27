package udatabase

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"

	// nolint:blank-imports // it is necessary to run the SQL migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs SQL migrations found in the folder specified by DBConfig.MigrationsDir
func RunMigrations(db *sql.DB, cfg *DBConfig) error {
	config := migratePostgres.Config{
		SchemaName: cfg.SchemaName,
	}
	driver, err := migratePostgres.WithInstance(db, &config)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsDir),
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
