package auth

import (
	"net/http"

	"github.com/Zigl3ur/go-app/internal/helper"
)

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handler for login endpoint
func (a *authService) loginHandler(w http.ResponseWriter, r *http.Request) {

	if isMethodAllowed := helper.MethodsAllowed(w, r, "POST"); !isMethodAllowed {
		return
	}

	// get body
	var body loginBody
	payload := helper.ReadBody(w, r, &body)

	// error if failed to parse body
	if payload == nil || body.Email == "" || body.Password == "" {
		helper.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to parse body, check fields"})
		return
	}

	// create session in database
	// no need to check rows since its checked in the func
	token, err := a.CreateSession(body.Email, body.Password)

	// response accordingly to database response
	switch {
	case err != nil:
		helper.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		cookie := http.Cookie{
			Name:     a.Config.CookieName,
			Value:    token,
			Path:     "/",
			MaxAge:   int(a.Config.SessionExpirity.Seconds()),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, &cookie)
		helper.JsonResponse(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

// handler for register endpoint
func (a *authService) registerHandler(w http.ResponseWriter, r *http.Request) {

	if isMethodAllowed := helper.MethodsAllowed(w, r, "POST"); !isMethodAllowed {
		return
	}

	// get body
	var body registerBody
	payload := helper.ReadBody(w, r, &body)

	// error if failed to parse body
	if payload == nil || body.Username == "" || body.Email == "" || body.Password == "" {
		helper.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "Failed to parse body, check fields"})
		return
	}

	err := a.CreateUser(body.Username, body.Email, body.Password)

	switch {
	case err != nil:
		helper.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		helper.JsonResponse(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

// handler for session endpoint
func (a *authService) getSession(w http.ResponseWriter, r *http.Request) {

	if isMethodAllowed := helper.MethodsAllowed(w, r, "GET"); !isMethodAllowed {
		return
	}

	token, err := r.Cookie(a.Config.CookieName)

	if err != nil {
		helper.JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "session cookie is missing"})
		return
	}

	session, user, err := a.CheckSession(token.Value)

	switch {
	case err != nil:
		helper.JsonResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	default:
		helper.JsonResponse(w, http.StatusOK, map[string]any{"session": &session, "user": &user})
	}

}

// handler for logout endpoint
func (a *authService) logoutHandler(w http.ResponseWriter, r *http.Request) {

	if isMethodAllowed := helper.MethodsAllowed(w, r, "GET"); !isMethodAllowed {
		return
	}

	token, err := r.Cookie(a.Config.CookieName)

	if err != nil {
		return
	}

	a.DeleteSession(token.Value)

	cookie := http.Cookie{
		Name:     a.Config.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	helper.JsonResponse(w, http.StatusOK, map[string]string{"status": "success"})
}
