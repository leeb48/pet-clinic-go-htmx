package app

import (
	"net/http"
	"runtime/debug"
)

func (app *App) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	data := app.NewTemplateData(r)
	app.Render(w, r, http.StatusInternalServerError, "server-error.html", data)
}

func (app *App) ClientError(w http.ResponseWriter, status int) {
	statusText := http.StatusText(status)

	app.Logger.Error("clientError", statusText, status)
	http.Error(w, statusText, status)
}

func (app *App) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}
