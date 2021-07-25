package render

import (
	"bytes"
	"fmt"
	"github.com/anthony-halim/booking-webapp/internal/config"
	"github.com/anthony-halim/booking-webapp/internal/models"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

// Map of functions that can be used in the template
var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from app config
		tc = app.TemplateCache
	} else {
		// Rebuild from cache
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	// Need to use make to initialise make, if not it's just a pointer to a nil map
	templateCache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	layoutMatches, err := filepath.Glob("./templates/*.layout.tmpl")
	if err != nil {
		return templateCache, err
	}
	hasLayoutMatches := len(layoutMatches) > 0

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		if hasLayoutMatches {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = ts
	}
	return templateCache, nil
}