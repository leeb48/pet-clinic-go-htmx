package handlers

import (
	"net/http"

	"pet-clinic.bonglee.com/internal/app"
)

type ErrorHandler struct {
	*app.App
}

func NewErrorHandler(app *app.App) *ErrorHandler {
	return &ErrorHandler{
		app,
	}
}

func (handler *ErrorHandler) serverError(w http.ResponseWriter, r *http.Request) {
	data := handler.NewTemplateData(r)
	handler.Render(w, r, http.StatusInternalServerError, "server-error.html", data)
}
