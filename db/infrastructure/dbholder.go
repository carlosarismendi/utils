package infrastructure

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/carlosarismendi/utils/shared/infrastructure/dotenv"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host                 string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port                 string `env:"POSTGRES_PORT" envDefault:"5432"`
	User                 string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password             string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	DatabaseName         string `env:"POSTGRES_DATABASE" envDefault:"postgres"`
	SchemaName           string `env:"POSTGRES_SCHEMA" envDefault:"public"`
	MigrationsDir        string `env:"POSTGRES_MIGRATIONS_DIR" envDefault:"./migrations"`
	RunMigrationsOnReset bool   `env:"POSTGRES_RUN_MIGRATIONS" envDefault:"false"`
}

func NewDBConfigFromEnv() *DBConfig {
	dotenv.Load()

	cfg := DBConfig{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}

func (c *DBConfig) checkValuesProvidedAndSetDefaults() {
	def := NewDBConfigFromEnv()

	if c.Host == "" {
		c.Host = def.Host
	}

	if c.Port == "" {
		c.Port = def.Port
	}

	if c.User == "" {
		c.User = def.User
	}

	if c.Password == "" {
		c.Password = def.Password
	}

	if c.DatabaseName == "" {
		c.DatabaseName = def.DatabaseName
	}

	if c.SchemaName == "" {
		c.SchemaName = def.SchemaName
	}

	if c.MigrationsDir == "" {
		c.MigrationsDir = def.MigrationsDir
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
	config *DBConfig
	db     *gorm.DB
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
		config: config,
		db:     db,
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
		SchemaName: d.config.SchemaName,
	}
	driver, err := migratePostgres.WithInstance(db, &config)
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s", d.config.MigrationsDir),
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

	err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", d.config.SchemaName)).Error
	if err != nil {
		panic(err)
	}

	d.createSchema()
	d.setSearchPath()

	if d.config.RunMigrationsOnReset {
		d.RunMigrations()
	}
}

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
