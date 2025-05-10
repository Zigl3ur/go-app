package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"time"
	"unicode/utf8"

	"github.com/Zigl3ur/go-app/internal/helper"
	"github.com/Zigl3ur/go-app/internal/store"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	errServerError         = errors.New("server errror")
	errInvalidEmail        = errors.New("email is invalid")
	errInvalidUsername     = errors.New("username is too short/long (min 3, max 20)")
	errInvalidCredentials  = errors.New("invalid credentials")
	errNoAccountFound      = errors.New("no account found")
	errAccountAlreadyExist = errors.New("account already exist")
	errCreatingSession     = errors.New("failed to create session")
	errDeletingSession     = errors.New("failed to delete session")
	errSessionNotFound     = errors.New("session not found")
	errSessionInvalid      = errors.New("session invalid")
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
func generateToken(length uint8) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errServerError
	}
	return hex.EncodeToString(b), nil
}

// create a user in the database
func (a *authService) createUser(username, email, password string) error {

	// check username length
	if utf8.RuneCountInString(username) <= 3 || utf8.RuneCountInString(username) >= 20 {
		return errInvalidUsername
	}

	// check email
	if _, err := mail.ParseAddress(email); err != nil {
		return errInvalidEmail
	}

	// check password
	if err := helper.CheckPassword(password); err != nil {
		return err
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
		if result.Error.Error() == "UNIQUE constraint failed: users.email" {
			return errAccountAlreadyExist
		}
		return errServerError
	}

	return nil
}

// update a user in the database
func (a *authService) updateUser(sessionToken, username, password string) error {

	// check password
	if err := helper.CheckPassword(password); err != nil {
		return err
	}

	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &store.User{
		Username:  username,
		Password:  string(hashPassword),
		UpdatedAt: time.Now(),
	}

	// check session to get userId from it
	_, userSession, err := a.getSession(sessionToken)
	if err != nil {
		return err
	}

	if result := a.Conn.Model(&store.User{}).Where(&store.User{Id: userSession.Id}).Updates(user); result.Error != nil {
		return errInvalidCredentials
	}
	return nil
}

// delete a user in the database
func (a *authService) deleteUser(sessionToken string) error {

	// check session to get userId from it
	_, user, err := a.getSession(sessionToken)
	if err != nil {
		return err
	}

	if result := a.Conn.Delete(&user); result.Error != nil {
		return errServerError
	}

	return nil
}

// create a session for given user credentials
func (a *authService) createSession(email, password string) (string, error) {

	var user store.User
	// retrieve user data
	if result := a.Conn.Select("id", "password").Where(&store.User{Email: email}).First(&user); result.Error != nil {
		return "", errNoAccountFound
	}

	// check password
	if isEqual := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); isEqual != nil {
		return "", errInvalidCredentials
	}

	token, err := generateToken(32)

	if err != nil {
		return "", errServerError
	}

	sessionTime := time.Now().Add(a.Config.SessionExpirity)

	session := &store.Session{
		Token:    token,
		UserId:   user.Id,
		ExpireAt: sessionTime,
	}

	// create the session
	if result := a.Conn.Create(&session); result.Error != nil {
		return "", errCreatingSession
	}

	return token, nil
}

// check if the given token is a valid session,
func (a *authService) getSession(token string) (*store.Session, *store.User, error) {

	var session store.Session
	var user store.User

	// TODO: join query maybe
	if result := a.Conn.Where(&store.Session{Token: token}).First(&session); result.Error != nil {
		fmt.Println(result.Error)
		return nil, nil, errSessionNotFound
	}

	if result := a.Conn.Where(&store.User{Id: session.UserId}).First(&user); result.Error != nil {
		return nil, nil, errNoAccountFound
	}

	// check expirity of session
	if isValid := time.Now().After(session.ExpireAt); isValid {
		return nil, nil, errSessionInvalid
	}

	return &session, &user, nil
}

// delete a session in the database,
func (a *authService) deleteSession(token string) error {

	var session store.Session

	if token == "" {
		return nil
	}

	if result := a.Conn.Where(&store.Session{Token: token}).Delete(&session); result.Error != nil {
		return errDeletingSession
	}

	return nil
}
