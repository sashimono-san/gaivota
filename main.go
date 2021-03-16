package main

import (
	"context"
	"fmt"
	"gaivota/handlers"
	"gaivota/internal/router"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func index(rw http.ResponseWriter, req *http.Request) {
	log.Println("Got request to /")
	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		// http.Error replaces the next couple of lines
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("Bad request."))
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}
	log.Printf("Data %s\n", data)
	fmt.Fprintf(rw, "Hello, %s!\n", data)
}

func main() {
	logger := log.New(os.Stdout, "Gaivota-api", log.LstdFlags)
	healthcheck := handlers.NewHealthCheck(logger)
	positions := handlers.NewPositions(logger)

	mux := router.New("/")
	mux.Get("/ping", healthcheck)
	mux.Get("/positions", http.HandlerFunc(positions.Get))
	mux.Post("/positions", http.HandlerFunc(positions.Add))

	// https://golang.org/pkg/net/http/#Server
	server := &http.Server{
		Addr:         ":9090",
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Println("Starting server")
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// https://golang.org/pkg/os/signal/#Notify
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// https://golang.org/ref/spec#Receive_operator
	sig := <-sigChan
	logger.Printf("Received terminate %s signal, gracefully shutting down.", sig)

	// https://pkg.go.dev/context
	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)
}
