package app

import "net/http"

func (app *App) Render(w http.ResponseWriter, r *http.Request, status int, page string, data TemplateData) {
	ts := app.TemplateCache[page]

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		panic(err)
	}
}
