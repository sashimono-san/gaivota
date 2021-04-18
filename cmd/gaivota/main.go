package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/leoschet/gaivota"
	"github.com/leoschet/gaivota/internal/config"
	"github.com/leoschet/gaivota/log"
	"github.com/leoschet/gaivota/mux"
	"github.com/leoschet/gaivota/postgres"
)

func main() {
	logger := log.New("Gaivota-api")

	settings, err := config.ReadFile("../config.json")
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, err)
	}

	if settings.Port == 0 {
		panic("Missing mandatory environment variable PORT")
	}

	db, err := postgres.Connect(context.Background(), settings.DatabaseConnString)
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, err)
	}
	defer db.Close()

	pgClient := db.NewPostgresClient()

	app := mux.New("")
	app.InitRouter(pgClient, []gaivota.HealthChecker{db}, logger)

	// https://golang.org/pkg/net/http/#Server
	server := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%v", settings.Port),
		Handler:      app.Router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Log(gaivota.LogLevelInfo, "Starting server")
		err := server.ListenAndServe()
		if err != nil {
			logger.Log(gaivota.LogLevelFatal, err)
		}
	}()

	// https://golang.org/pkg/os/signal/#Notify
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// https://golang.org/ref/spec#Receive_operator
	sig := <-sigChan
	logger.Log(gaivota.LogLevelInfo, "Received terminate %s signal, gracefully shutting down.", sig)

	// https://pkg.go.dev/context
	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeoutContext)
}
