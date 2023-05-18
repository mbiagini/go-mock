package main

import (
	"context"
	"fmt"
	"go-mock/model"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/mbiagini/go-server-utils/gslog"

	"github.com/go-chi/chi/v5"
)

func main() {

	err := LoadConfiguration("./resources/config.json")
	if err != nil {
		fmt.Println("Error found in app configuration")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = model.ValidateEndpoints(Conf.Endpoints)
	if err != nil {
		gslog.Server(err.Error())
		os.Exit(1)
	}

	quit := make(chan interface{})
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv := startServer(httpServerExitDone, quit)

	<-quit
	gslog.Server("main: shutting down.")
	if err := srv.Shutdown(context.TODO()); err != nil {
		gslog.Server(fmt.Sprintf("error while trying to shutdown: %s", err.Error()))
		os.Exit(1)
	}

	// wait for goroutine started in startServer() to stop
	httpServerExitDone.Wait()

	gslog.Server("main: done shuting down. Exiting.")
}

func startServer(wg *sync.WaitGroup, quit chan interface{}) *http.Server {

	r := chi.NewRouter()

	// Initialize routes.
	gslog.Server("Initializing routes")
	Routes(r, quit)

	// Print routes.
	gslog.Server("Printing routes")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		gslog.Server(fmt.Sprintf("[%s]: '%s'", method, route))
		return nil
	})

	// Start server.
	gslog.Server("Starting server")
	srv := &http.Server{Addr: Conf.Ip + ":" + strconv.Itoa(Conf.Port), Handler: r}
	go func() {
		defer wg.Done() // let main know we are done cleaning up.

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error.
			gslog.Server(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()

	return srv
}