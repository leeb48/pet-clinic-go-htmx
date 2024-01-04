package main

import "net/http"

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts := app.templateCache[page]

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		panic(err)
	}
}

// todo: validation error is not showing when status is written in the response
