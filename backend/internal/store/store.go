package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// connect to DB,
// return a *gormDB and an error
func Connect(sqlitedb string, gormConfig gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(sqlitedb), &gormConfig)

	if err != nil {
		return nil, err
	}

	return db, nil
}
