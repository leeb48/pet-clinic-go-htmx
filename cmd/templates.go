package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"pet-clinic.bonglee.com/ui"
)

type Alert struct {
	Msg     string
	MsgType string
}
type templateData struct {
	Alert    Alert
	FlashMsg string
	Form     any
}

const (
	FLASH_MSG = "FLASH_MSG"
	PRIMARY   = "primary"
	DANGER    = "danger"
	SUCCESS   = "success"
	WARNING   = "warning"
	INFO      = "info"
)

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		FlashMsg: app.session.PopString(r.Context(), FLASH_MSG),
	}
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
