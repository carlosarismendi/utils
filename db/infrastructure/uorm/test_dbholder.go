package uorm

import (
	"fmt"

	"github.com/carlosarismendi/utils/db/infrastructure"
)

type TestDBHolder struct {
	*DBHolder
}

func NewTestDBHolder(schemaName string) *TestDBHolder {
	cfg := infrastructure.NewDBConfigFromEnv()
	cfg.SchemaName = schemaName
	return &TestDBHolder{
		DBHolder: NewDBHolder(cfg),
	}
}

func (d *TestDBHolder) Reset() {
	db := d.db

	err := db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", d.config.SchemaName)).Error
	if err != nil {
		panic(err)
	}

	sdb, err := db.DB()
	if err != nil {
		panic(err)
	}
	d.config.CreateSchema(sdb)
	d.config.SetSearchPath(sdb)
	_ = d.RunMigrations()
}
