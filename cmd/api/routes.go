package main

import "github.com/julienschmidt/httprouter"

func (app *application) route() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", app.Index)
	router.GET("/hello/:name", app.Hello)
	return router
}
