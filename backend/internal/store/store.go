package store

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	Conn *gorm.DB
}

func (s *Store) Connect(sqlitedb string, gormConfig gorm.Config) error {
	db, err := gorm.Open(sqlite.Open(sqlitedb), &gormConfig)

	if err != nil {
		return err
	}

	s.Conn = db
	return nil
}
