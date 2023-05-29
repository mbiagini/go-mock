package main

import (
	"go-mock/apierrors"
	"go-mock/controller"
	"go-mock/db"
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
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		gsrender.WriteJSON(w, http.StatusNotFound, apierrors.New(apierrors.OPERATION_NOT_DEFINED))
	})
}

func addServerRoutes(r *chi.Mux, quit chan interface{}) {

	gslog.Server("Initializing server routes")

	r.Route("/go-mock/v1", func(r chi.Router) {

		r.Post("/restart", func(w http.ResponseWriter, r *http.Request) {
			close(quit)
		})

		r.Route("/endpoints", func(r chi.Router) {
			r.Get("/", controller.GetEndpoints)
			r.Get("/{id}", controller.GetEndpointById)
			r.Post("/", controller.PostEndpoint)
			r.Put("/{id}", controller.UpdateEndpoint)
			r.Delete("/{id}", controller.DeleteEndpoint)
		})

		r.Route("/files", func(r chi.Router) {
			r.Get("/", controller.GetAllFiles)
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				controller.PostFiles(w, r, Conf.UploadMaxSize)
			})
		})

	})

}

func addUserRoutes(r *chi.Mux) {

	gslog.Server("Initializing user defined routes")

	// All routes will begin with the server basepath.
	r.Route(Conf.Basepath, func(r chi.Router) {

		for _, e := range db.DB.FindAll() {
			r.MethodFunc(string(e.Method), e.Path, addRoute(e))
		}

	})

}

// addRoute creates a new Handler for the given Endpoint.
func addRoute(e model.Endpoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.HandleRequest(w, r, e)
	}
}