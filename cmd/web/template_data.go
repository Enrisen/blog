package main

import (
	"github.com/Enrisen/blog/internal/data"
)

type TemplateData struct {
	Title      string
	HeaderText string
	FormErrors map[string]string
	FormData   map[string]string
	Posts      []*data.Post // Add Posts field for blog
	CSRFToken  string       // Add CSRF token field
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Information Technology Blog",
		HeaderText: "Welcome to my Information Technology Blog",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
		Posts:      []*data.Post{}, // Initialize Posts
	}
}
