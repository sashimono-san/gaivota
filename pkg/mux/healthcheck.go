package mux

import (
	"log"
	"net/http"

	"github.com/leoschet/gaivota"
)

type HealthCheck struct {
	logger       *log.Logger
	dependencies []gaivota.HealthChecker
}

func NewHealthCheck(logger *log.Logger, dependencies []gaivota.HealthChecker) *HealthCheck {
	return &HealthCheck{
		logger,
		dependencies,
	}
}

func (hc *HealthCheck) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, dep := range hc.dependencies {
		msg, err := dep.Ping()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(msg))
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("pong"))
}
