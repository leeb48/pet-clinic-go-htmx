package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()

	dynamic := alice.New(app.session.LoadAndSave, app.logRequest)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/owner", dynamic.ThenFunc(app.ownerList))

	router.Handler(http.MethodGet, "/owner/create", dynamic.ThenFunc(app.ownerCreate))
	router.Handler(http.MethodPost, "/owner/create", dynamic.ThenFunc(app.ownerCreatePost))

	router.Handler(http.MethodPost, "/pet/add-pet-type", dynamic.ThenFunc(app.newPetTypePost))

	router.GET("/admin", app.adminPage)

	return router
}
