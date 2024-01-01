package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"pet-clinic.bonglee.com/internal/models"
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

func (app *application) ownerCreatePost(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

	var owner models.Owner
	err := json.NewDecoder(r.Body).Decode(&owner)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	fmt.Printf("%#v", owner)

	w.Header().Add("HX-Redirect", "/owner/create")
}
