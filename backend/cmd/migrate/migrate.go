package main

import (
	"log"

	"github.com/Zigl3ur/go-app/internal/store"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&store.User{},
		&store.Session{},
	)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	store := &store.Store{}

	err := store.Connect("db.sql", gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = migrate(store.Conn)

	if err != nil {
		log.Fatal(err)
	}
}
