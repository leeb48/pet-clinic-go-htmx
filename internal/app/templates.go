package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"pet-clinic.bonglee.com/ui"
)

type Alert struct {
	Msg     string
	MsgType string
}
type TemplateData struct {
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

func (app *App) NewTemplateData(r *http.Request) TemplateData {
	return TemplateData{
		FlashMsg: app.Session.PopString(r.Context(), FLASH_MSG),
	}
}

func phoneNumber(phone string) string {
	formatPhone := fmt.Sprintf("(%s) %s-%s", phone[0:3], phone[3:6], phone[6:])

	return formatPhone
}

func birthDate(date time.Time) string {
	return date.Format("01/02/2006")
}

var functions = template.FuncMap{
	"phoneNumber": phoneNumber,
	"birthdate":   birthDate,
}

func CreateTemplateCache() (map[string]*template.Template, error) {
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

		ts := template.Must(template.
			New(name).
			Funcs(functions).
			ParseFS(ui.Files, patterns...))

		cache[name] = ts
	}

	return cache, nil
}
