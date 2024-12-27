package main

import (
	"net/http"
)

var routes = []*route{}

type route struct {
	path        string
	handlerFunc func(w http.ResponseWriter, r *http.Request)
	middlewares []func(http.Handler) http.Handler
}

// util function for new route creation
func registerRoute(path string, handlerFunc http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	routes = append(routes, &route{
		path:        path,
		handlerFunc: handlerFunc,
		middlewares: middlewares,
	})
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// register routes

	// static (add middleware sanitizer)
	registerRoute("/img/", app.img, app.staticResourcePathValidatorMiddleware, app.staticResourceIsModifiedMiddleware)
	registerRoute("/css/", app.css, app.staticResourcePathValidatorMiddleware, app.staticResourceIsModifiedMiddleware)
	registerRoute("/js/", app.js, app.staticResourcePathValidatorMiddleware, app.staticResourceIsModifiedMiddleware)

	// endpoints
	registerRoute("/", app.home)
	registerRoute("/auth", app.auth, app.authMiddleware)

	// set handler for every route and middlewares if any
	for _, route := range routes {
		handler := http.HandlerFunc(route.handlerFunc)
		mux.Handle(route.path, app.chainMiddleware(handler, route.middlewares...))
	}

	// handler with global middlewares applied
	return app.chainMiddleware(mux, app.loggerMiddleware)
}
