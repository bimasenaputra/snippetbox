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

func add(i1, i2 int) int {
	return i1 + i2
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"add": add,
}

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	HasNext bool
	HasPrev bool
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

	fragments, err := filepath.Glob("./ui/html/fragments/*.html")
	if err != nil {
		return nil, err
	}

	for _, fragment := range fragments {

		name := filepath.Base(fragment)

		ts, err := template.New(fragment).Funcs(functions).ParseFiles(fragment)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	
	return cache, nil
}