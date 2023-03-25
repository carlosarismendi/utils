package infrastructure

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/carlosarismendi/utils/shared/dotenv"
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

// NewDBConfigFromEnv returns a *DBConfig initialized by env variables
func NewDBConfigFromEnv() *DBConfig {
	dotenv.Load()

	cfg := DBConfig{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
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

	if c.MigrationsDir == "" {
		c.MigrationsDir = "./migrations"
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
