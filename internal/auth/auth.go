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

var (
	errInvalidEmail       = errors.New("email is invalid")
	errInvalidUsername    = errors.New("username is too short/long (min 3, max 20)")
	errInvalidPassword    = errors.New("password is invalid, it must be min 8 chars long, contain 1 special char, 1 upper and lower char and 1 digit")
	errInvalidCredentials = errors.New("invalid credentials")
	errGenerateToken      = errors.New("failed to generate token")
	errNoAccountFound     = errors.New("no account found")
	errCreatingSession    = errors.New("failed to create session")
	errDeletingSession    = errors.New("failed to delete session")
	errSessionNotFound    = errors.New("session not found")
	errSessionInvalid     = errors.New("session invalid")
)

type authService struct {
	Conn   *gorm.DB
	Config *config
}

type config struct {
	Endpoint        string
	CookieName      string
	SessionExpirity time.Duration
}

func NewAuthService(db *gorm.DB, endpoint, cookiename string, expirity time.Duration) *authService {

	newConfig := &config{CookieName: cookiename, Endpoint: endpoint, SessionExpirity: expirity}
	newService := &authService{Conn: db, Config: newConfig}

	if endpoint == "" {
		newConfig.Endpoint = "/api/auth"
	}

	if cookiename == "" {
		newConfig.CookieName = "session"
	}

	return newService
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
func (a *authService) CreateUser(username, email, password string) (int64, error) {

	// check username length
	if utf8.RuneCountInString(username) <= 3 || utf8.RuneCountInString(username) >= 20 {
		return 0, errInvalidUsername
	}

	// check email
	if _, err := mail.ParseAddress(email); err != nil {
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
		return result.RowsAffected, result.Error
	}

	return result.RowsAffected, nil
}

// delete a user in the database
// return rowsAffected and an error
func (a *authService) DeleteUser(username, email, password string) (int64, error) {

	var user store.User
	result := a.Conn.Select("password").Where(&store.User{Username: username, Email: email}).Find(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	// check password
	if isEqual := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); isEqual != nil {
		return 0, errInvalidCredentials
	}

	result = a.Conn.Where(&user).Delete(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// create a session for given user credentials
// return the token and an error
func (a *authService) CreateSession(email, password string) (string, error) {

	var user store.User
	// retrieve user data
	result := a.Conn.Select("id, password").Where("email = ?", email).First(&user)

	if result.Error != nil {
		return "", errNoAccountFound
	}

	// check password
	if isEqual := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); isEqual != nil {
		return "", errInvalidCredentials
	}

	token, err := generateToken(32)

	if err != nil {
		return "", errGenerateToken
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
		return "", errCreatingSession
	}

	return token, nil
}

// check if the given token is a valid session,
// return true if valid, false for expired and an error
func (a *authService) CheckSession(token string) (bool, error) {

	var session store.Session
	result := a.Conn.Select("token").Where(&store.Session{Token: token}).Find(&session)

	if result.Error != nil {
		return false, result.Error

	} else if result.RowsAffected == 0 {
		return false, errSessionNotFound
	}

	// check expirity of session
	if isValid := !session.ExpiresAt.Equal(time.Now()); !isValid {
		return false, errSessionInvalid
	}

	return true, nil
}

// delete a session in the database,
// return rowsAffected and an error
func (a *authService) DeleteSession(token string) (int64, error) {

	var session store.Session
	result := a.Conn.Where(&store.Session{Token: token}).Delete(&session)

	if result.Error != nil {
		return 0, errDeletingSession
	}

	return result.RowsAffected, nil
}
