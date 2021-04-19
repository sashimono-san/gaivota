package mux

import (
	"net/http"
	"reflect"

	"github.com/leoschet/gaivota"
)

func InitHealthCheckRouter(mux *Mux, dependencies []gaivota.HealthChecker, logger gaivota.Logger) {
	healthcheck := &HealthCheck{
		logger,
		dependencies,
	}

	mux.Router.Get("/ping", healthcheck)
}

type HealthCheck struct {
	logger       gaivota.Logger
	dependencies []gaivota.HealthChecker
}

func (hc *HealthCheck) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	hc.logger.Log(gaivota.LogLevelInfo, "Handle ping endpoint")

	for _, dep := range hc.dependencies {
		msg, err := dep.Ping()
		if err != nil {
			hc.logger.Log(gaivota.LogLevelInfo, "Error while pinging %s: %v\n", reflect.TypeOf(dep).Name(), err)
			http.Error(rw, msg, http.StatusInternalServerError)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("pong"))
}
