package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) ownerList(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "owner-list.html", data)
}

func (app *application) ownerCreate(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "owner-create.html", data)
}
