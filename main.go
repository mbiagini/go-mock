package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	quit := make(chan interface{})
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv := startServer(httpServerExitDone, quit)

	<-quit
	log.Println("main: shutting down")
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err)
	}

	// wait for goroutine started in startServer() to stop
	httpServerExitDone.Wait()

	log.Println("main: done. exiting")
}

func startServer(wg *sync.WaitGroup, quit chan interface{}) *http.Server {

	r := chi.NewRouter()

	// Initialize routes.
	log.Println("Initializing routes")
	r.Route(Conf.Basepath, func(r chi.Router) {
		r.MethodFunc("POST", "/shutdown", func(w http.ResponseWriter, r *http.Request) {
			close(quit)
		})
		for _, ep := range Conf.Endpoints {
			r.MethodFunc(ep.Method, ep.Path, addRoute(ep))
		}
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
	})

	// Print routes.
	log.Println("Printing routes")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[%s]: '%s'\n", method, route)
		return nil
	})

	// Start server.
	log.Println("Starting server")
	srv := &http.Server{Addr: Conf.Ip + ":" + strconv.Itoa(Conf.Port), Handler: r}
	go func() {
		defer wg.Done() // let main know we are done cleaning up.

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error.
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

func addRoute(ep Endpoint) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, ep)
	}
}

func respond(w http.ResponseWriter, r *http.Request, ep Endpoint) {

	var response Response
	
	if (!ep.HasDiscriminator) {
		response = findResponseById(ep.Responses, ep.DefaultResponseId)
	} else {
		d := ep.Discriminator
		var paramValue string
		switch (strings.ToLower(d.Location)) {
		case "path":
			paramValue = chi.URLParam(r, d.Parameter)
		case "query":
			paramValue = r.URL.Query().Get(d.Parameter)
		case "header":
			paramValue = r.Header.Get(d.Parameter)
		default:
			log.Println("Parameter location found not one of path, query or header after config validation")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		found := false
		for _, c := range d.Conditions {
			if (paramValue == c.Value) {
				response = findResponseById(ep.Responses, c.ResponseId)
				found = true
				break
			}
		}
		if (!found) {
			log.Println("No match found in conditions for parameter", paramValue, ". Using default response")
			response = findResponseById(ep.Responses, ep.DefaultResponseId)
		}
	}

	if (response.Delay > 0) {
		time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	}

	if (response.ContentType != "") {
		w.Header().Set("Content-Type", response.ContentType)
	}
	w.WriteHeader(response.Code)
	if (response.BodyFilename != "") {
		w.Write(readFile(response.BodyFilename))
	}
}

func findResponseById(rs []Response, id int) Response {
	for _, r := range rs {
		if (id == r.Id) {
			return r
		}
	}
	panic("No response found for given default response id after config validation")
}

func readFile(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if (err != nil) {
		log.Println("Could not read file", filename)
	}
	return buf
}