package usql

import (
	"github.com/carlosarismendi/utils/udatabase"

	// nolint:blank-imports // it is necessary to run the SQL migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
)

type DBHolder struct {
	config *udatabase.DBConfig
	db     *sqlx.DB
}

// Returns a *DBHolder initialized with the provided config.
// In case the *DBConfig object has zero values, those will
// be filled with default values.
func NewDBHolder(config *udatabase.DBConfig) *DBHolder {
	config.SetEmptyValuesToDefaults()

	conn := config.GetConnectionString()
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		panic(err)
	}

	dbHolder := &DBHolder{
		config: config,
		db:     db,
	}

	dbHolder.config.CreateSchema(dbHolder.db.DB)
	dbHolder.config.SetSearchPath(dbHolder.db.DB)

	return dbHolder
}

// RunMigrations runs SQL migrations found in the folder specified by DBConfig.MigrationsDir
func (d *DBHolder) RunMigrations() error {
	return udatabase.RunMigrations(d.db.DB, d.config)
}

// GetDBInstance returns the inner database object *sqlx.DB provided by sqlx.
// More on sqlx here: https://github.com/jmoiron/sqlx
func (d *DBHolder) GetDBInstance() *sqlx.DB {
	return d.db
}
