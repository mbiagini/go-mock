package main

import (
	"errors"
	"go-mock/controller"
	"go-mock/model"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mbiagini/go-server-utils/gslog"
	"github.com/mbiagini/go-server-utils/gsrender"
)

// Routes initializes all the routes. These include the default ones (which we will
// call server routes) and the ones loaded from the configuration file (those would
// be the user routes).
//
// Receives a channel to notice when /shutdown endpoint is called (server route).
func Routes(r *chi.Mux, quit chan interface{}) {
	addUserRoutes(r)
	addServerRoutes(r, quit)
}

func addServerRoutes(r *chi.Mux, quit chan interface{}) {

	gslog.Server("Initializing server routes")

	// Enables the programatically shutdown of the server.
	r.MethodFunc("POST", "/shutdown", func(w http.ResponseWriter, r *http.Request) {
		close(quit)
	})

}

func addUserRoutes(r *chi.Mux) {

	gslog.Server("Initializing user defined routes")

	// All routes will begin with the server basepath.
	r.Route(Conf.Basepath, func(r chi.Router) {

		for _, e := range Conf.Endpoints {
			r.MethodFunc(string(e.Method), e.Path, addRoute(e))
		}

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			gsrender.WriteJSON(w, http.StatusNotFound, model.ErrorFrom(errors.New("operation not found")))
		})

	})

}

// addRoute creates a new Handler for the given Endpoint.
func addRoute(e model.Endpoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.HandleRequest(w, r, e)
	}
}