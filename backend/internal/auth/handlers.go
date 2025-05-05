package auth

import (
	"net/http"

	"github.com/Zigl3ur/go-app/internal/helper"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handler for Login Logic
func (a *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {

	isMethodAllowed := helper.MethodsAllowed(w, r, "POST")

	if !isMethodAllowed {
		return
	}

	// get body
	var body LoginBody
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
