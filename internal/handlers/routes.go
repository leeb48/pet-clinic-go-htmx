package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/ui"
)

func Routes(app *app.App) *httprouter.Router {

	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	dynamic := alice.New(app.Session.LoadAndSave, app.LogRequest)

	ownerHandler := NewOwnerHandler(app)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(ownerHandler.home))
	router.Handler(http.MethodGet, "/owner", dynamic.ThenFunc(ownerHandler.list))
	router.Handler(http.MethodGet, "/owner/create", dynamic.ThenFunc(ownerHandler.createPage))
	router.Handler(http.MethodPost, "/owner/create", dynamic.ThenFunc(ownerHandler.createPost))
	router.Handler(http.MethodGet, "/owner/edit/:id", dynamic.ThenFunc(ownerHandler.editPage))
	router.Handler(http.MethodPut, "/owner/edit/:id", dynamic.ThenFunc(ownerHandler.editPut))
	router.Handler(http.MethodGet, "/owner/detail/:id", dynamic.ThenFunc(ownerHandler.detail))
	router.Handler(http.MethodDelete, "/owner/:id", dynamic.ThenFunc(ownerHandler.remove))

	petHandler := NewPetHandler(app)
	router.Handler(http.MethodGet, "/admin", dynamic.ThenFunc(petHandler.adminPage))
	router.Handler(http.MethodPost, "/pet/add-pet-type", dynamic.ThenFunc(petHandler.newPetTypePost))
	router.Handler(http.MethodPost, "/pet/search", dynamic.ThenFunc(petHandler.getPetsByNameAndDob))

	vetHandler := NewVetHandler(app)
	router.Handler(http.MethodGet, "/vet", dynamic.ThenFunc(vetHandler.list))
	router.Handler(http.MethodGet, "/vet/create", dynamic.ThenFunc(vetHandler.createPage))
	router.Handler(http.MethodPost, "/vet/create", dynamic.ThenFunc(vetHandler.createPost))
	router.Handler(http.MethodGet, "/vet/detail/:id", dynamic.ThenFunc(vetHandler.detail))
	router.Handler(http.MethodGet, "/vet/edit/:id", dynamic.ThenFunc(vetHandler.edit))
	router.Handler(http.MethodPut, "/vet/edit/:id", dynamic.ThenFunc(vetHandler.editPost))
	router.Handler(http.MethodDelete, "/vet/:id", dynamic.ThenFunc(vetHandler.remove))
	router.Handler(http.MethodGet, "/vet/search", dynamic.ThenFunc(vetHandler.getByLastName))

	visitHandler := NewVisitHandler(app)
	router.Handler(http.MethodGet, "/visit/detail/:id", dynamic.ThenFunc(visitHandler.visitDetail))
	router.Handler(http.MethodGet, "/visit/edit/:id", dynamic.ThenFunc(visitHandler.editVisitPage))
	router.Handler(http.MethodPost, "/visit/create", dynamic.ThenFunc(visitHandler.createVisitPost))
	router.Handler(http.MethodGet, "/visit/vetId/:id", dynamic.ThenFunc(visitHandler.getVisitCalendarByVetId))
	router.Handler(http.MethodDelete, "/visit/:visitId/vetId/:vetId", dynamic.ThenFunc(visitHandler.removeVisit))

	return router
}
