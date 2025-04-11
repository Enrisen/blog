package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.blogPage)

	// Blog post routes
	mux.HandleFunc("GET /blog/create", app.blogCreateForm)
	mux.HandleFunc("POST /blog/create", app.blogCreateSubmit)

	// User registration routes
	mux.HandleFunc("GET /user/register", app.userRegisterForm)
	mux.HandleFunc("POST /user/register", app.userRegisterSubmit)

	return app.loggingMiddleware(mux)
}
