package auth

import "net/http"

func (a *authService) Router() {
	http.HandleFunc(a.Config.Endpoint+"/login", a.LoginHandler)
	http.HandleFunc(a.Config.Endpoint+"/register", a.RegisterHandler)
	http.HandleFunc(a.Config.Endpoint+"/session", a.GetSession)
}
