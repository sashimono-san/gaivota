package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/leoschet/gaivota"
	"github.com/leoschet/gaivota/internal/config"
	"github.com/leoschet/gaivota/log"
	"github.com/leoschet/gaivota/mux"
	"github.com/leoschet/gaivota/postgres"
)

func main() {
	logger := log.New("Gaivota-api - ")

	rootPath, err := os.Getwd()
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, "Error getting root path: %v", err)
	}

	settings, err := config.ReadFile(path.Join(rootPath, "/config.json"))
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, "Error while reading config file: %v", err)
	}

	if settings.Port == 0 {
		panic("Missing mandatory environment variable PORT")
	}

	db, err := postgres.Connect(context.Background(), settings.DatabaseConnString)
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, "Error while connecting to Postgres: %v", err)
	}
	defer db.Close()

	pgClient := db.NewPostgresClient()

	app := mux.New("/")
	app.InitRouter(pgClient, []gaivota.HealthChecker{db}, logger)

	addr := fmt.Sprintf("0.0.0.0:%v", settings.Port)
	// https://golang.org/pkg/net/http/#Server
	server := &http.Server{
		Addr:         addr,
		Handler:      app.Router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		logger.Log(gaivota.LogLevelInfo, "Starting server on %s", addr)
		err := server.ListenAndServe()
		if err != nil {
			logger.Log(gaivota.LogLevelFatal, "Error while starting server: %v", err)
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
