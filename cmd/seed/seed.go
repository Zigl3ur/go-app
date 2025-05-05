package main

import (
	"log"

	"github.com/Zigl3ur/go-app/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {

	paswordHash, _ := bcrypt.GenerateFromPassword([]byte("AzertyuioP1234!"), 10)

	user := &store.User{
		Username: "john doe",
		Email:    "john@doe.com",
		Password: string(paswordHash),
	}

	db, err := store.Connect("db.sql", gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	if result := db.Create(user); result.Error != nil {
		log.Fatal("Failed to create user:", result.Error)
	}

	log.Println("Database seeded")
}
