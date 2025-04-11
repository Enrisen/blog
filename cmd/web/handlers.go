package main

import (
	"net/http"
	"strings"

	"github.com/Enrisen/blog/internal/validator"
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

// blogCreateForm displays the blog post creation form
func (app *application) blogCreateForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Create Post | TechSphere"
	data.HeaderText = "Create a New Blog Post"

	err := app.render(w, http.StatusOK, "blog_create.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render blog creation page", "template", "blog_create.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// blogCreateSubmit processes the blog post creation form submission
func (app *application) blogCreateSubmit(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse blog creation form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract form values
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	categoriesStr := r.PostForm.Get("categories")

	// Process categories (split by comma and trim spaces)
	var categories []string
	if categoriesStr != "" {
		for _, category := range strings.Split(categoriesStr, ",") {
			trimmedCategory := strings.TrimSpace(category)
			if trimmedCategory != "" {
				categories = append(categories, trimmedCategory)
			}
		}
	}

	// Initialize form data and errors
	data := NewTemplateData()
	data.Title = "Create Post | TechSphere"
	data.HeaderText = "Create a New Blog Post"
	data.FormData = map[string]string{
		"title":      title,
		"content":    content,
		"categories": categoriesStr,
	}

	// Validate form inputs using the validator
	v := validator.NewValidator()
	v.ValidateBlogPost(title, content, categories)

	// If there are validation errors, re-render the form
	if !v.ValidData() {
		data.FormErrors = v.Errors
		err = app.render(w, http.StatusUnprocessableEntity, "blog_create.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render blog creation page with errors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Use dummy author ID for now (Enrisen Tzib)
	var authorID int64 = 1

	// Create the blog post
	_, err = app.blog.CreatePost(authorID, title, content, categories)
	if err != nil {
		app.logger.Error("failed to create blog post", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Redirect to the home page instead of showing success page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// userRegisterForm displays the user registration form
func (app *application) userRegisterForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Register | TechSphere"
	data.HeaderText = "Create an Account"

	err := app.render(w, http.StatusOK, "user_register.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render user registration page", "template", "user_register.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// userRegisterSubmit processes the user registration form submission
func (app *application) userRegisterSubmit(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse registration form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract form values
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirmPassword := r.PostForm.Get("confirm_password")

	// Initialize form data and errors
	data := NewTemplateData()
	data.Title = "Register | TechSphere"
	data.HeaderText = "Create an Account"
	data.FormData = map[string]string{
		"name":  name,
		"email": email,
	}

	// Validate form inputs using the validator
	v := validator.NewValidator()
	v.ValidateUserRegistration(name, email, password, confirmPassword)

	// If there are validation errors, re-render the form
	if !v.ValidData() {
		data.FormErrors = v.Errors
		err = app.render(w, http.StatusUnprocessableEntity, "user_register.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render user registration page with errors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Register the user (password will be hashed in the model)
	_, err = app.users.RegisterUser(name, email, password, confirmPassword)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			// Create a new validator for the duplicate email error
			v := validator.NewValidator()
			v.AddError("email", "Email address is already in use")
			data.FormErrors = v.Errors
			err = app.render(w, http.StatusUnprocessableEntity, "user_register.tmpl", data)
			if err != nil {
				app.logger.Error("failed to render user registration page with duplicate email error", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		app.logger.Error("failed to register user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the success page
	data = NewTemplateData()
	data.Title = "Registration Successful | TechSphere"
	data.HeaderText = "Registration Successful"

	err = app.render(w, http.StatusOK, "user_register_success.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render registration success page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
