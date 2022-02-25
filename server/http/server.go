package http

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	hc "questions-rest-api/interfaceAdapters/http"
	"strconv"
	"time"
)

// RegisterRoutes registers the http server routes
func RegisterRoutes(r *mux.Router, c hc.Controller) {
	r.HandleFunc("/question", c.Add).Methods("POST")
	r.HandleFunc("/question/{id:[0-9]+}", c.Update).Methods("PUT")
	r.HandleFunc("/question/{id:[0-9]+}", c.Delete).Methods("DELETE")
	r.HandleFunc("/questions", c.GetAll).Methods("GET")

	// create Redoc configuration
	ops := middleware.RedocOpts{
		SpecURL: "/swagger.yaml",
	}

	// add swagger documentation routes
	sh := middleware.Redoc(ops, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
}

// StartServer starts a new http server that listens on the given port
func StartServer(r *mux.Router, port int) {
	s := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start server on a different goroutine
	go func() {
		log.Println("Starting server on port " + strconv.Itoa(port))

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(fmt.Sprintf("unable to start http server: %s", err.Error()))
		}

	}()

	// create a signal channel that will be notified for Interrupt and Kill signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// wait for a signal
	sig := <-sigChan
	log.Println("Received terminate, graceful shutdown", sig)

	// create context with timeout, the server will wait 30 seconds for all connections to finish
	tc, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := s.Shutdown(tc); err != nil {
		log.Fatalln(fmt.Sprintf("error shuting down server: %s", err.Error()))
	}

}
