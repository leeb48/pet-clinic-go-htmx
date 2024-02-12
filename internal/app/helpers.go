package app

import (
	"bytes"
	"fmt"
	"net/http"
)

func (app *App) Render(w http.ResponseWriter, r *http.Request, status int, page string, data TemplateData) {
	ts, ok := app.TemplateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, r, err)
		return
	}

	// use a buffer to render the template to check for any errors
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
