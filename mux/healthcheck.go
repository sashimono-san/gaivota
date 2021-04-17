package mux

import (
	"log"
	"net/http"

	"github.com/leoschet/gaivota"
)

func NewHealthCheck(logger *log.Logger, dependencies []gaivota.HealthChecker) *HealthCheck {
	return &HealthCheck{
		logger,
		dependencies,
	}
}

type HealthCheck struct {
	logger       *log.Logger
	dependencies []gaivota.HealthChecker
}

func (hc *HealthCheck) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, dep := range hc.dependencies {
		msg, err := dep.Ping()
		if err != nil {
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("pong"))
}
