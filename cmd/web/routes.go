package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.blogPage)

	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeaderMiddleware,
	)

	// Blog post routes
	mux.HandleFunc("GET /blog/view/{id}", app.blogView)
	mux.HandleFunc("GET /blog/create", app.blogCreateForm)
	mux.HandleFunc("POST /blog/create", app.blogCreateSubmit)
	mux.HandleFunc("GET /blog/edit/{id}", app.blogEditForm)
	mux.HandleFunc("POST /blog/edit/{id}", app.blogEditSubmit)
	mux.HandleFunc("POST /blog/delete/{id}", app.blogDelete)

	// User registration routes
	mux.HandleFunc("GET /user/register", app.userRegisterForm)
	mux.HandleFunc("POST /user/register", app.userRegisterSubmit)

	dynamicMiddleware := application.New(app.session.Enable)

	return app.loggingMiddleware(mux)
}
