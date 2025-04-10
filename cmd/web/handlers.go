package main

import (
	"net/http"
)

func (app *application) blogPage(w http.ResponseWriter, r *http.Request) {
	// Get all blog posts
	posts, err := app.blog.GetAll()
	if err != nil {
		app.logger.Error("failed to get blog posts", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "TechSphere | Information Technology Blog"
	data.HeaderText = "Latest Technology Articles"
	data.Posts = posts

	err = app.render(w, http.StatusOK, "blog.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render blog page", "template", "blog.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
