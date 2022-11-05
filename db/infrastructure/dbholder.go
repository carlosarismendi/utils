package infrastructure

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBHolder struct {
	schemaName string
	db         *gorm.DB
}

func NewDBHolderFromInstance(db *gorm.DB) *DBHolder {	
	rows, err := db.Debug().Raw("SHOW search_path").Rows()
	if err != nil {
		panic(err)
	}

	var currentSchema string
	rows.Next() 
	err = rows.Scan(&currentSchema)
	if err != nil {
		panic(err)
	}

	return &DBHolder{
		schemaName: currentSchema,
		db: db,
	}
}

func NewDBHolder(schemaName string) *DBHolder {
	conn := fmt.Sprintf("host=localhost user=postgres password=postgres dbname=postgres search_path=%s port=5432 sslmode=disable TimeZone=UTC", schemaName)
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	dbHolder := &DBHolder{
		schemaName: schemaName,
		db:         db,
	}

	dbHolder.Reset()
	return dbHolder
}

func (d *DBHolder) Reset() {
	db := d.db

	err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", d.schemaName)).Error
	if err != nil {
		panic(err)
	}

	d.createSchema()
	d.setSearchPath()
	d.runMigrations()
}

func (d *DBHolder) createSchema() {
	err := d.db.Exec(fmt.Sprintf("CREATE SCHEMA %s;", d.schemaName)).Error
	if err != nil {
		panic(err)
	}
}

func (d *DBHolder) setSearchPath() {
	d.db = d.db.Exec(fmt.Sprintf("SET search_path TO %s;", d.schemaName))
	if d.db.Error != nil {
		panic(d.db.Error)
	}

}
func (d *DBHolder) runMigrations() {
	db, err := d.db.DB()
	if err != nil {
		panic(err)
	}

	config := migratePostgres.Config{
		SchemaName: d.schemaName,
	}
	driver, err := migratePostgres.WithInstance(db, &config)
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:/home/microservices/dev/ddd-hexa/migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		panic(err)
	}
}

func (d *DBHolder) GetDBInstance() *gorm.DB {
	return d.db
}
