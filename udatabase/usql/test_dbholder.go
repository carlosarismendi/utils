package usql

import (
	"fmt"

	"github.com/carlosarismendi/utils/udatabase"
)

type TestDBHolder struct {
	*DBHolder
}

func NewTestDBHolder(schemaName string) *TestDBHolder {
	cfg := udatabase.NewDBConfigFromEnv()
	cfg.SchemaName = schemaName
	return &TestDBHolder{
		DBHolder: NewDBHolder(cfg),
	}
}

func (d *TestDBHolder) Reset() {
	_, err := d.db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", d.config.SchemaName))
	if err != nil {
		panic(err)
	}

	d.config.CreateSchema(d.db.DB)
	d.config.SetSearchPath(d.db.DB)
	_ = d.RunMigrations()
}
