package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zigl3ur/go-app/internal/auth"
	"github.com/Zigl3ur/go-app/internal/store"
	"gorm.io/gorm"
)

func main() {

	db := &store.Store{}
	db.Connect("db.sql", gorm.Config{})

	auth := &auth.AuthService{
		Conn: db.Conn,
	}

	rowsAffected, err := auth.CreateUser("eden", "e@ee.com", "1234")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(rowsAffected)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Print("a") })

	log.Fatal(http.ListenAndServe(":8000", nil))
}
