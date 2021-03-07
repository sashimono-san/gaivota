package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type HealthCheck struct {
	logger *log.Logger
}

func NewHealthCheck(logger *log.Logger) *HealthCheck {
	return &HealthCheck{logger}
}

func (healthChecker *HealthCheck) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	message := "pong"
	healthChecker.logger.Println(message)
	fmt.Fprint(rw, message)
}
