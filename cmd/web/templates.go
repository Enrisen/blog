package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {
	// initialize an empty map
	cache := map[string]*template.Template{}
	// let's find all the templates. We will get the
	// full path for the files in a slice
	pages, err := filepath.Glob("./ui/html/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)
		// parse the page
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// make an entry into the cache
		cache[fileName] = ts
	}

	return cache, nil
}
