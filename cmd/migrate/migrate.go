package main

import (
	"log"

	"github.com/Zigl3ur/go-app/internal/store"
	"gorm.io/gorm"
)

// push schema to database,
// return an error
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

	db, err := store.Connect("db.sql", gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	err = migrate(db)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Schema Migration Done")
}
