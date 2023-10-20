package db

import (
	"path/filepath"

	"github.com/indig0fox/a3go/assemblyfinder"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var configuredDbPath string = filepath.Join(
	filepath.Dir(assemblyfinder.GetModulePath()),
	"stats.db",
)

func Client() *gorm.DB {
	if db == nil {
		err := Connect(configuredDbPath)
		if err != nil {
			panic(err)
		}
	}
	return db
}

func Connect(dbPath string) error {
	var err error
	configuredDbPath = dbPath
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
