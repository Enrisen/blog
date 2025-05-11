package main

import (
	"net/http"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
		app.logger.Info("Request processed")
	})
	return fn

}

// requireAuthentication middleware checks if a user is authenticated
// If not, it redirects to the login page
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		if app.session.GetInt(r, "authenticatedUserID") == 0 {
			// User is not authenticated, redirect to login page
			app.session.Put(r, "flash", "You must be logged in to access this page")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Set the "Cache-Control: no-store" header to prevent caching of protected pages
		// This helps prevent users from using the back button to access protected pages after logout
		w.Header().Add("Cache-Control", "no-store")

		// User is authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}

// secureHeaders adds security headers to the response
func (app *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}
