package infrastructure

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"

	// nolint:blank-imports // it is necessary to run the SQL migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBHolder struct {
	config *DBConfig
	db     *gorm.DB
}

// Returns a *DBHolder initialized with the provided config.
// In case the *DBConfig object has zero values, those will
// be filled with default values.
func NewDBHolder(config *DBConfig) *DBHolder {
	config.checkValuesProvidedAndSetDefaults()

	conn := config.getConnectionString()
	pg := postgres.Open(conn)
	db, err := gorm.Open(pg, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	dbHolder := &DBHolder{
		config: config,
		db:     db,
	}

	dbHolder.createSchema()
	dbHolder.setSearchPath()

	return dbHolder
}

// RunMigrations runs SQL migrations found in the folder specified by DBConfig.MigrationsDir
func (d *DBHolder) RunMigrations() error {
	db, err := d.db.DB()
	if err != nil {
		return err
	}

	config := migratePostgres.Config{
		SchemaName: d.config.SchemaName,
	}
	driver, err := migratePostgres.WithInstance(db, &config)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", d.config.MigrationsDir),
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

// GetDBInstance returns the inner database object *gorm.DB provided by GORM.
// More on GORM here: https://gorm.io/
func (d *DBHolder) GetDBInstance() *gorm.DB {
	return d.db
}

func (d *DBHolder) createSchema() {
	err := d.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", d.config.SchemaName)).Error
	if err != nil {
		panic(err)
	}
}

func (d *DBHolder) setSearchPath() {
	d.db = d.db.Exec(fmt.Sprintf("SET search_path TO %s;", d.config.SchemaName))
	if d.db.Error != nil {
		panic(d.db.Error)
	}
}
