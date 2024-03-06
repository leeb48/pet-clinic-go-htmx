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
	router.Handler(http.MethodDelete, "/owner/:id", dynamic.ThenFunc(ownerHandler.ownerRemove))

	petHandler := NewPetHandler(app)
	router.Handler(http.MethodGet, "/admin", dynamic.ThenFunc(petHandler.adminPage))
	router.Handler(http.MethodPost, "/pet/add-pet-type", dynamic.ThenFunc(petHandler.newPetTypePost))
	router.Handler(http.MethodPost, "/pet/search", dynamic.ThenFunc(petHandler.getPetsByNameAndDob))

	vetHandler := NewVetHandler(app)
	router.Handler(http.MethodGet, "/vet", dynamic.ThenFunc(vetHandler.vetList))
	router.Handler(http.MethodGet, "/vet/create", dynamic.ThenFunc(vetHandler.vetCreate))
	router.Handler(http.MethodPost, "/vet/create", dynamic.ThenFunc(vetHandler.vetCreatePost))
	router.Handler(http.MethodGet, "/vet/detail/:id", dynamic.ThenFunc(vetHandler.vetDetail))
	router.Handler(http.MethodGet, "/vet/edit/:id", dynamic.ThenFunc(vetHandler.vetEdit))
	router.Handler(http.MethodPut, "/vet/edit/:id", dynamic.ThenFunc(vetHandler.vetEditPut))
	router.Handler(http.MethodDelete, "/vet/:id", dynamic.ThenFunc(vetHandler.vetRemove))
	router.Handler(http.MethodPost, "/vet/visit", dynamic.ThenFunc(vetHandler.vetCreateVisitPost))

	return router
}
