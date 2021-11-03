package database

import (
	"github.com/pkg/errors"
	"github.com/sunreaver/antman/v2/db"
	"gorm.io/gorm"
)

var (
	dbs     map[string]*db.Databases
	configs map[string]db.Config
)

func InitDateBase(config map[string]db.Config) (err error) {
	dbs = make(map[string]*db.Databases, len(config))
	for tag, c := range config {
		d, err := db.MakeDB(c, &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			return errors.Wrap(err, "init database")
		}

		dbs[tag] = d
	}
	configs = config

	return nil
}

func AFIREMaster() *gorm.DB {
	return dbs["afire"].Master()
}

func AFIRESlave() *gorm.DB {
	return dbs["afire"].Slave()
}
