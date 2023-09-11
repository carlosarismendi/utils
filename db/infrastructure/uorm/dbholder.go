package uorm

import (
	"context"

	"github.com/carlosarismendi/utils/db/infrastructure"

	// nolint:blank-imports // it is necessary to run the SQL migrations.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBHolder struct {
	config *infrastructure.DBConfig
	db     *gorm.DB
}

// Returns a *DBHolder initialized with the provided config.
// In case the *DBConfig object has zero values, those will
// be filled with default values.
func NewDBHolder(config *infrastructure.DBConfig) *DBHolder {
	config.SetEmptyValuesToDefaults()

	conn := config.GetConnectionString()
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

	sdb, err := db.DB()
	if err != nil {
		panic(err)
	}
	config.CreateSchema(sdb)
	config.SetSearchPath(sdb)

	return dbHolder
}

// RunMigrations runs SQL migrations found in the folder specified by DBConfig.MigrationsDir
func (d *DBHolder) RunMigrations() error {
	sdb, err := d.db.DB()
	if err != nil {
		panic(err)
	}
	return infrastructure.RunMigrations(sdb, d.config)
}

// GetDBInstance returns the inner database object *gorm.DB provided by GORM.
// More on GORM here: https://gorm.io/
func (d *DBHolder) GetDBInstance(ctx context.Context) *gorm.DB {
	txFromCtx := ctx.Value(ctxk(transactionName))
	if txFromCtx == nil {
		txFromCtx = d.db.WithContext(ctx)
	}
	return txFromCtx.(*gorm.DB)
}
