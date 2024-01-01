package main

import "github.com/julienschmidt/httprouter"

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()
	router.GET("/", app.home)

	router.GET("/owner", app.ownerList)
	router.GET("/owner/create", app.ownerCreate)
	router.POST("/owner/create", app.ownerCreatePost)

	return router
}
