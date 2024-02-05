package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"pet-clinic.bonglee.com/internal/app"
)

func Routes(app *app.App) *httprouter.Router {

	router := httprouter.New()

	dynamic := alice.New(app.Session.LoadAndSave, app.LogRequest)

	ownerHandler := NewOwnerHandler(app)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(ownerHandler.home))
	router.Handler(http.MethodGet, "/owner", dynamic.ThenFunc(ownerHandler.ownerList))
	router.Handler(http.MethodGet, "/owner/create", dynamic.ThenFunc(ownerHandler.ownerCreate))
	router.Handler(http.MethodPost, "/owner/create", dynamic.ThenFunc(ownerHandler.ownerCreatePost))
	router.Handler(http.MethodGet, "/owner/edit/:id", dynamic.ThenFunc(ownerHandler.ownerEdit))
	router.Handler(http.MethodPut, "/owner/edit/:id", dynamic.ThenFunc(ownerHandler.ownerEditPut))
	router.Handler(http.MethodGet, "/owner/detail/:id", dynamic.ThenFunc(ownerHandler.ownerDetail))
	router.Handler(http.MethodDelete, "/owner/:id", dynamic.ThenFunc(ownerHandler.RemoveOwner))

	petHandler := NewPetHandler(app)
	router.Handler(http.MethodGet, "/admin", dynamic.ThenFunc(petHandler.adminPage))
	router.Handler(http.MethodPost, "/pet/add-pet-type", dynamic.ThenFunc(petHandler.newPetTypePost))

	return router
}
