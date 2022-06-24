package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.bimasenaputra/internal/models"
)

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
}

func newTemplateCache() (map[string]*template.Template, error) {	
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	cache := map[string]*template.Template{}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(page).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts 
	}

	return cache, nil
}