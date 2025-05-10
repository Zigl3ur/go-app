package auth

import "net/http"

func (a *authService) Router() {
	http.HandleFunc(a.Config.Endpoint+"/login", a.loginHandler)
	http.HandleFunc(a.Config.Endpoint+"/register", a.registerHandler)
	http.HandleFunc(a.Config.Endpoint+"/session", a.getSession)
	http.HandleFunc(a.Config.Endpoint+"/logout", a.logoutHandler)
	http.HandleFunc(a.Config.Endpoint+"/user", a.updateUserHandler)
}
