package auth

import "net/http"

func (a *authService) Router() {

	// auth not needed routes
	http.HandleFunc(a.Config.Endpoint+"/login", a.loginHandler)
	http.HandleFunc(a.Config.Endpoint+"/register", a.registerHandler)

	// auth needed handlers
	sessionHandler := http.HandlerFunc(a.getSessionHandler)
	updateUserHandler := http.HandlerFunc(a.updateUserHandler)
	logoutHandler := http.HandlerFunc(a.logoutHandler)

	// auth needed routes
	http.Handle(a.Config.Endpoint+"/session", a.authMiddleware(sessionHandler))
	http.Handle(a.Config.Endpoint+"/user", a.authMiddleware(updateUserHandler))
	http.Handle(a.Config.Endpoint+"/logout", a.authMiddleware(logoutHandler))
}
