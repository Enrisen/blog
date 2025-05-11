package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *TemplateData) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	// Add default data to the template data
	data = app.addDefaultData(r, data)

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.logger.Error("template does not exist", "template", page, "error", err)
		return err
	}
	err := ts.Execute(buf, data)
	if err != nil {
		err = fmt.Errorf("failed to render template %s: %w", page, err)
		app.logger.Error("failed to render template", "template", page, "error", err)
		return err
	}

	w.WriteHeader(status)

	_, err = buf.WriteTo(w)
	if err != nil {
		err = fmt.Errorf("failed to write template to response: %w", err)
		app.logger.Error("failed to write template to response", "error", err)
		return err
	}

	return nil
}

// addDefaultData adds default data to the template data
func (app *application) addDefaultData(r *http.Request, data *TemplateData) *TemplateData {
	// If no data was passed, initialize a new empty TemplateData
	if data == nil {
		data = NewTemplateData()
	}

	// Add flash message to the template data if one exists
	data.Flash = app.session.PopString(r, "flash")

	// Add authentication status to the template data
	userID := app.session.GetInt(r, "authenticatedUserID")
	if userID > 0 {
		data.IsAuthenticated = true
		data.UserName = app.session.GetString(r, "userName")

		// Log authentication status for debugging
		app.logger.Info("User is authenticated",
			"userID", userID,
			"userName", data.UserName,
			"path", r.URL.Path)
	} else {
		// Log that user is not authenticated
		app.logger.Info("User is not authenticated",
			"path", r.URL.Path,
			"session_data", app.session.Exists(r, "authenticatedUserID"))
	}

	return data
}

// serverError logs the detailed error message and sends a generic 500 Internal Server Error response to the client.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(), "url", r.URL.Path, "method", r.Method)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and corresponding description to the client.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound sends a 404 Not Found response to the client.
func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, http.StatusNotFound)
}
