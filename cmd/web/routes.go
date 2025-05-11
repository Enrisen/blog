package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Public routes
	mux.HandleFunc("GET /{$}", app.blogPage)
	mux.HandleFunc("GET /blog/view/{id}", app.blogView)
	mux.HandleFunc("GET /user/register", app.userRegisterForm)
	mux.HandleFunc("POST /user/register", app.userRegisterSubmit)
	mux.HandleFunc("GET /user/login", app.userLoginForm)
	mux.HandleFunc("POST /user/login", app.userLoginSubmit)
	mux.HandleFunc("GET /user/logout", app.userLogout)

	// Protected routes
	mux.Handle("GET /blog/create", app.requireAuthentication(http.HandlerFunc(app.blogCreateForm)))
	mux.Handle("POST /blog/create", app.requireAuthentication(http.HandlerFunc(app.blogCreateSubmit)))
	mux.Handle("GET /blog/edit/{id}", app.requireAuthentication(http.HandlerFunc(app.blogEditForm)))
	mux.Handle("POST /blog/edit/{id}", app.requireAuthentication(http.HandlerFunc(app.blogEditSubmit)))
	mux.Handle("POST /blog/delete/{id}", app.requireAuthentication(http.HandlerFunc(app.blogDelete)))

	// Apply session middleware, then noSurf, and finally logging middleware
	return app.loggingMiddleware(app.session.Enable(noSurf(mux)))
}
