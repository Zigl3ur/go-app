package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Zigl3ur/go-app/internal/auth"
	"github.com/Zigl3ur/go-app/internal/store"
	"gorm.io/gorm"
)

func main() {

	db, err := store.Connect("db.sql", gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	auth := &auth.AuthService{
		Conn: db,
		Config: &auth.Config{
			CookieName:      "session",
			SessionExpirity: time.Duration(24 * time.Hour),
		},
	}

	log.Fatal(http.ListenAndServe(":8000", nil))
}
