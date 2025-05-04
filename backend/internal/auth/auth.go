package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/mail"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/Zigl3ur/go-app/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const cost uint8 = 15

var (
	errInvalidEmail       = errors.New("email is invalid")
	errInvalidUsername    = errors.New("username is too short/long (min 3, max 20)")
	errInvalidPassword    = errors.New("password is invalid, it must be min 8 chars long, contain 1 special char, 1 upper and lower char and 1 digit")
	errGenerateToken      = errors.New("failed to generate token")
	errInvalidCredentials = errors.New("invalid credentials")
	errNoAccountFound     = errors.New("no account found")
	errCreatingSession    = errors.New("failed to create session")
)

type AuthService struct {
	Conn   *gorm.DB
	Config *Config
}

type Config struct {
	CookieName      string
	SessionExpirity time.Duration
}

// generate a hexa decimal token from given length,
// return the generated token and an error
func generateToken(length uint8) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errGenerateToken
	}
	return hex.EncodeToString(b), nil
}

// create a user in the database,
// return rowsAffected and an error
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

	// password regex check
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasDigit, _ := regexp.MatchString(`\d`, password)
	hasSpecial, _ := regexp.MatchString(`[^\da-zA-Z]`, password)
	hasLength := utf8.RuneCountInString(password) >= 8

	if !(hasLower && hasUpper && hasDigit && hasSpecial && hasLength) {
		return 0, errInvalidPassword
	}

	// hash password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	// create user struct
	user := &store.User{
		Username: username,
		Email:    email,
		Password: string(passwordHash),
	}

	// insert user into db
	result := a.Conn.Create(&user)

	if result.Error != nil {
		return result.RowsAffected, err
	}

	return result.RowsAffected, nil
}

// create a session for given user credentials
// return rowsAffected and an error
func (a *AuthService) CreateSession(email, password string) (int64, error) {

	var user store.User
	// retrieve user data
	result := a.Conn.Select("id, password").Where("email = ?", email).First(&user)

	if result.Error != nil {
		return 0, errNoAccountFound
	}

	// check password
	isEqual := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if isEqual != nil {
		return 0, errInvalidCredentials
	}

	token, err := generateToken(32)

	if err != nil {
		return 0, errGenerateToken
	}

	sessionTime := time.Now().Add(a.Config.SessionExpirity)

	session := &store.Session{
		Token:     token,
		UserId:    user.Id,
		ExpiresAt: sessionTime,
	}

	// create the session
	result = a.Conn.Create(&session)

	if result.Error != nil {
		return 0, errCreatingSession
	}

	return result.RowsAffected, nil
}
