package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"pet-clinic.bonglee.com/ui"
)

type templateData struct {
	Form any
}

func createTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/**/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/layout/nav.html",
			page,
		}

		ts := template.Must(template.ParseFS(ui.Files, patterns...))

		cache[name] = ts
	}

	return cache, nil
}
