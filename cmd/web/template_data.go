package main

import (
	"github.com/Enrisen/blog/internal/data"
)

type TemplateData struct {
	Title              string
	HeaderText         string
	FormErrors         map[string]string
	FormData           map[string]string
	Posts              []*data.Post     // Add Posts field for blog
	Post               *data.Post       // Single post for blog view
	RecentPosts        []*data.Post     // Recent posts for sidebar
	RelatedPosts       []*data.Post     // Related posts for blog view
	Categories         []*data.Category // All available categories for sidebar/forms
	SelectedCategories map[string]bool  // Map of category names that are selected for a post (for edit form)
	CSRFToken          string           // Add CSRF token field
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:              "Information Technology Blog",
		HeaderText:         "Welcome to my Information Technology Blog",
		FormErrors:         map[string]string{},
		FormData:           map[string]string{},
		Posts:              []*data.Post{},     // Initialize Posts
		Post:               nil,                // Initialize Post as nil
		RecentPosts:        []*data.Post{},     // Initialize RecentPosts
		RelatedPosts:       []*data.Post{},     // Initialize RelatedPosts
		Categories:         []*data.Category{}, // Initialize Categories
		SelectedCategories: map[string]bool{},  // Initialize SelectedCategories
	}
}
