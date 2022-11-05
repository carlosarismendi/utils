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

type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SchemaName   string
}

func (c *DBConfig) checkValuesProvidedAndSetDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == "" {
		c.Port = "5432"
	}

	if c.User == "" {
		c.User = "postgres"
	}

	if c.Password == "" {
		c.Password = "postgres"
	}

	if c.DatabaseName == "" {
		c.DatabaseName = "postgres"
	}

	if c.SchemaName == "" {
		c.SchemaName = "public"
	}
}

func (c *DBConfig) getConnectionString() string {
	host := fmt.Sprintf("host=%s", c.Host)
	port := fmt.Sprintf("port=%s", c.Port)
	user := fmt.Sprintf("user=%s", c.User)
	pass := fmt.Sprintf("password=%s", c.Password)
	dbname := fmt.Sprintf("dbname=%s", c.DatabaseName)
	search_path := fmt.Sprintf("search_path=%s", c.SchemaName)

	conn := fmt.Sprintf("%s %s %s %s %s %s sslmode=disable TimeZone=UTC", host, port, user, pass, dbname, search_path)
	return conn
}

type DBHolder struct {
	schemaName string
	db         *gorm.DB
}

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
		schemaName: config.SchemaName,
		db:         db,
	}

	dbHolder.createSchema()
	dbHolder.setSearchPath()

	return dbHolder
}

func (d *DBHolder) RunMigrations() {
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

func (d *DBHolder) Reset() {
	db := d.db

	err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", d.schemaName)).Error
	if err != nil {
		panic(err)
	}

	d.createSchema()
	d.setSearchPath()
	d.RunMigrations()
}

func (d *DBHolder) GetDBInstance() *gorm.DB {
	return d.db
}

func (d *DBHolder) createSchema() {
	err := d.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", d.schemaName)).Error
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
