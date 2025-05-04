package auth

import (
	"errors"
	"net/mail"
	"unicode/utf8"

	"github.com/Zigl3ur/go-app/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	errInvalidEmail    = errors.New("email is invalid")
	errInvalidUsername = errors.New("username is too short/long (min 3, max 20)")
)

type AuthService struct {
	Conn *gorm.DB
}

// function to create a user in the database, return rowsAffected and error
func (a *AuthService) CreateUser(username, email, password string) (int64, error) {

	// check username length
	if utf8.RuneCountInString(username) <= 3 || utf8.RuneCountInString(username) >= 20 {
		return 0, errInvalidUsername
	}

	_, err := mail.ParseAddress(email)

	// check email
	if err != nil {
		return 0, errInvalidEmail
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return 0, err
	}

	// create user struct
	user := &store.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// insert user into db
	result := a.Conn.Select("Username", "Email", "Password").Create(&user)

	if result.Error != nil {
		return result.RowsAffected, err
	}

	return result.RowsAffected, nil
}
