package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"pet-clinic.bonglee.com/internal/constants/alertConstants"
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

func (app *App) NewTemplateData(r *http.Request) TemplateData {
	return TemplateData{
		FlashMsg: app.Session.PopString(r.Context(), alertConstants.FLASH_MSG),
	}
}

func phoneNumber(phone string) string {
	formatPhone := fmt.Sprintf("(%s) %s-%s", phone[0:3], phone[3:6], phone[6:])

	return formatPhone
}

func birthDate(date time.Time) string {
	return date.Format("01/02/2006")
}

func YYYYMMDD(date time.Time) string {
	return date.Format("2006-01-02")
}

func add(x, y int) int {
	return x + y
}

func toDateTime(dateTime time.Time) string {
	parsedTime := dateTime.Local().Format(time.DateTime)
	return parsedTime
}

var functions = template.FuncMap{
	"phoneNumber": phoneNumber,
	"birthdate":   birthDate,
	"YYYYMMDD":    YYYYMMDD,
	"add":         add,
	"toDateTime":  toDateTime,
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
