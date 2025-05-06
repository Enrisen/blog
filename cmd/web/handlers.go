package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	appdata "github.com/Enrisen/blog/internal/data"
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

	// Fetch all categories for the sidebar
	categories, err := app.blog.GetAllCategories()
	if err != nil {
		app.logger.Error("failed to get categories", "error", err)
		// Don't fail the whole page, just log the error
	}

	data := NewTemplateData()
	data.Title = "TechSphere | Information Technology Blog"
	data.HeaderText = "Latest Technology Articles"
	data.Categories = categories // Add categories to template data
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
	// Fetch all categories for the dropdown
	categories, err := app.blog.GetAllCategories()
	if err != nil {
		app.logger.Error("failed to get categories for create form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData()
	data.Title = "Create Post | TechSphere"
	data.HeaderText = "Create a New Blog Post"
	data.Categories = categories // Add categories to template data

	err = app.render(w, http.StatusOK, "blog_create.tmpl", data) // Use = instead of :=
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
	selectedCategory := r.PostForm.Get("category")
	newCategory := strings.TrimSpace(r.PostForm.Get("new_category"))

	// Determine which category to save
	categoriesToSave := []string{}

	if selectedCategory == "new" && newCategory != "" {
		// User selected "Add new category" and provided a name
		categoriesToSave = append(categoriesToSave, newCategory)
	} else if selectedCategory != "" && selectedCategory != "new" {
		// User selected an existing category
		categoriesToSave = append(categoriesToSave, selectedCategory)
	}
	// If neither condition is met, no category will be saved

	// Initialize form data and errors
	data := NewTemplateData()
	data.Title = "Create Post | TechSphere"
	data.HeaderText = "Create a New Blog Post"
	data.FormData = map[string]string{
		"title":        title,
		"content":      content,
		"new_category": newCategory,
		// We might need a way to pass selectedCategories back if needed for re-population
	}
	// TODO: Add selectedCategories back to FormData if needed for re-populating multi-select

	// Validate form inputs using the validator
	v := validator.NewValidator()
	// Pass the combined list for validation
	appdata.ValidateBlogPost(v, title, content, categoriesToSave)
	// TODO: Add specific validation for newCategory if desired (e.g., length)

	// If there are validation errors, re-render the form
	if !v.ValidData() {
		// Fetch categories again to populate the form dropdown on error
		categoriesForForm, catErr := app.blog.GetAllCategories()
		if catErr != nil {
			app.logger.Error("failed to get categories for create form re-render", "error", catErr)
			// If we can't get categories, render the form without them, but log the error
		}

		// Create a map of selected categories for the template
		data.SelectedCategories = make(map[string]bool)
		if selectedCategory != "" && selectedCategory != "new" {
			data.SelectedCategories[selectedCategory] = true
		}

		data.Categories = categoriesForForm
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
	_, err = app.blog.CreatePost(authorID, title, content, categoriesToSave) // Pass combined list
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
	appdata.ValidateUserRegistration(v, name, email, password, confirmPassword)

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

// blogView displays a single blog post
func (app *application) blogView(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Get the post from the database
	post, err := app.blog.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			app.logger.Error("failed to get blog post", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Get recent posts for the sidebar
	recentPosts, err := app.blog.GetAll()
	if err != nil {
		app.logger.Error("failed to get recent posts", "error", err)
		// Continue with the page even if we can't get recent posts
	}

	// Fetch all categories for the sidebar
	categories, err := app.blog.GetAllCategories()
	if err != nil {
		app.logger.Error("failed to get categories for view page", "error", err)
		// Don't fail the whole page, just log the error
	}

	// Prepare template data
	data := NewTemplateData()
	data.Title = post.Title + " | TechSphere"
	data.HeaderText = post.Title
	data.Categories = categories // Add categories to template data
	data.Post = post
	data.RecentPosts = recentPosts

	// Render the template
	err = app.render(w, http.StatusOK, "blog_view.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render blog view page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// blogEditForm displays the blog post edit form
func (app *application) blogEditForm(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Get the post from the database
	post, err := app.blog.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			app.logger.Error("failed to get blog post", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Fetch all categories for the dropdown
	allCategories, err := app.blog.GetAllCategories()
	if err != nil {
		app.logger.Error("failed to get categories for edit form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := NewTemplateData()
	data.Title = "Edit Post | TechSphere"
	data.HeaderText = "Edit Blog Post"
	data.Categories = allCategories // Add all categories for the dropdown
	data.Post = post

	// Create a map of the post's current categories for easy lookup in the template
	data.SelectedCategories = make(map[string]bool)

	// For the dropdown, we'll just select the first category if there is one
	if len(post.Categories) > 0 {
		data.SelectedCategories[post.Categories[0]] = true
	}

	// Pre-populate form data
	data.FormData = map[string]string{
		"title":        post.Title,
		"content":      post.Content,
		"new_category": "", // Initialize empty for the new category field
	}

	// Render the template
	err = app.render(w, http.StatusOK, "blog_edit.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render blog edit page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// blogEditSubmit processes the blog post edit form submission
func (app *application) blogEditSubmit(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Parse the form data
	err = r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse blog edit form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract form values
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	selectedCategory := r.PostForm.Get("category")
	newCategory := strings.TrimSpace(r.PostForm.Get("new_category"))

	// Determine which category to save
	categoriesToSave := []string{}

	if selectedCategory == "new" && newCategory != "" {
		// User selected "Add new category" and provided a name
		categoriesToSave = append(categoriesToSave, newCategory)
	} else if selectedCategory != "" && selectedCategory != "new" {
		// User selected an existing category
		categoriesToSave = append(categoriesToSave, selectedCategory)
	}
	// If neither condition is met, no category will be saved

	// Initialize form data and errors
	data := NewTemplateData()
	data.Title = "Edit Post | TechSphere"
	data.HeaderText = "Edit Blog Post"
	data.FormData = map[string]string{
		"title":        title,
		"content":      content,
		"new_category": newCategory,
		// TODO: Add selectedCategories back to FormData if needed for re-populating multi-select
	}

	// Validate form inputs using the validator
	v := validator.NewValidator()
	appdata.ValidateBlogPost(v, title, content, categoriesToSave) // Pass combined list
	// TODO: Add specific validation for newCategory if desired

	// If there are validation errors, re-render the form
	if !v.ValidData() {
		// Fetch categories again to populate the form dropdown on error
		categoriesForForm, catErr := app.blog.GetAllCategories()
		if catErr != nil {
			app.logger.Error("failed to get categories for edit form re-render", "error", catErr)
			// If we can't get categories, render the form without them, but log the error
		}
		// Need to pass the post data again as well
		data.Post = &appdata.Post{ID: id} // Pass a minimal post struct just for ID if needed, or fetch again if more data required by template
		data.Categories = categoriesForForm
		// Create a map of selected categories for the template
		data.SelectedCategories = make(map[string]bool)
		if selectedCategory != "" && selectedCategory != "new" {
			data.SelectedCategories[selectedCategory] = true
		}

		data.FormErrors = v.Errors
		err = app.render(w, http.StatusUnprocessableEntity, "blog_edit.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render blog edit page with errors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Update the blog post
	err = app.blog.UpdatePost(id, title, content, categoriesToSave) // Pass combined list
	if err != nil {
		app.logger.Error("failed to update blog post", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the updated blog post
	http.Redirect(w, r, "/blog/view/"+idStr, http.StatusSeeOther)
}

// blogDelete handles the deletion of a blog post
func (app *application) blogDelete(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Delete the post
	err = app.blog.DeletePost(id)
	if err != nil {
		app.logger.Error("failed to delete blog post", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the blog page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// userLoginForm displays the user login form
func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Login | TechSphere"
	data.HeaderText = "Log In to Your Account"

	err := app.render(w, http.StatusOK, "user_login.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render user login page", "template", "user_login.tmpl", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// userLoginSubmit processes the user login form submission
func (app *application) userLoginSubmit(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse login form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract form values
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// Initialize form data and errors
	data := NewTemplateData()
	data.Title = "Login | TechSphere"
	data.HeaderText = "Log In to Your Account"
	data.FormData = map[string]string{
		"email": email,
	}

	// Validate form inputs using the validator
	v := validator.NewValidator()
	appdata.ValidateLogin(v, email, password)

	// If there are validation errors, re-render the form
	if !v.ValidData() {
		data.FormErrors = v.Errors
		err = app.render(w, http.StatusUnprocessableEntity, "user_login.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render user login page with errors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Authenticate the user
	user, err := app.users.Authenticate(email, password)
	if err != nil {
		// Create a new validator for the authentication error
		v := validator.NewValidator()
		v.AddError("email", "Invalid email or password")
		data.FormErrors = v.Errors
		err = app.render(w, http.StatusUnprocessableEntity, "user_login.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render user login page with authentication error", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// TODO: Set up a session for the authenticated user
	app.logger.Info("user logged in successfully", "user_id", user.ID, "email", user.Email)

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
