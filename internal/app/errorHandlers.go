package app

import (
	"net/http"
	"runtime/debug"
)

type errorForm struct {
	ErrorMsg   string
	StatusText string
	StatusCode int
}

func (app *App) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.Logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	data := app.NewTemplateData(r)
	data.Form = &errorForm{
		ErrorMsg: err.Error(),
	}
	app.Render(w, r, http.StatusInternalServerError, "server-error.html", data)
}

func (app *App) ClientError(w http.ResponseWriter, r *http.Request, statusCode int) {
	statusText := http.StatusText(statusCode)

	app.Logger.Error("clientError", statusText, statusCode)

	data := app.NewTemplateData(r)
	data.Form = &errorForm{
		StatusText: statusText,
		StatusCode: statusCode,
	}

	app.Render(w, r, http.StatusInternalServerError, "client-error.html", data)
}
