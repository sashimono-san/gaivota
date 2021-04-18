package mux

import (
	"github.com/leoschet/gaivota"
	"github.com/leoschet/mux"
)

func New(prefix string) *Mux {
	return &Mux{
		Router: mux.NewRouter(prefix),
	}
}

type Mux struct {
	Router mux.Router
}

func (mux *Mux) InitRouter(client *gaivota.Client, dependencies []gaivota.HealthChecker, logger *gaivota.Logger) {
	InitHealthCheckRouter(mux, dependencies, logger)
	InitPortfolioRouter(mux, client.PortfolioStore, logger)
}

// TODO: Create initial Router here, get a Client as parameter
// 				Every handler sets its own routes
