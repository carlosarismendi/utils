package infrastructure

import (
	"fmt"
)

type TestDBHolder struct {
	*DBHolder
}

func NewTestDBHolder(schemaName string) *TestDBHolder {
	cfg := NewDBConfigFromEnv()
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

	d.createSchema()
	d.setSearchPath()

	if d.config.RunMigrationsOnReset {
		d.RunMigrations()
	}
}
