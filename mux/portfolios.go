package mux

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/leoschet/gaivota"
	"github.com/leoschet/mux"
)

func InitPortfolioRouter(mux *Mux, store gaivota.PortfolioStore, logger *gaivota.Logger) {
	portfolioHandler := &PortfolioHandler{
		logger:         logger,
		PortfolioStore: store,
	}

	router := mux.Router.NewSubrouter("/positions")

	router.Post("/", http.HandlerFunc(portfolioHandler.Add))
	router.Get("/:portfolioId", http.HandlerFunc(portfolioHandler.Get))
}

type PortfolioHandler struct {
	logger         gaivota.Logger
	PortfolioStore gaivota.PortfolioStore
}

func (handler *PortfolioHandler) Get(rw http.ResponseWriter, req *http.Request) {
	handler.logger.Log(gaivota.LogLevelInfo, "Handle GET Portfolio")

	params := mux.PathParams(req)
	portfolioId, err := strconv.Atoi(params["portfolioId"])

	if err != nil {
		http.Error(rw, "Portfolio ID must be an integer", http.StatusBadRequest)
		return
	}

	portfolio, err := handler.PortfolioStore.Get(context.Background(), portfolioId)

	if err != nil {
		http.Error(rw, "Error while getting Portfolio", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(portfolio)
}

func (handler *PortfolioHandler) Add(rw http.ResponseWriter, req *http.Request) {
	handler.logger.Log(gaivota.LogLevelInfo, "Handle POST Portfolio")

	var portfolio gaivota.Portfolio
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&portfolio)

	// TODO: Improve error handling as in: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
	if err != nil {
		handler.logger.Log(gaivota.LogLevelInfo, "Error while decoding POST /portfolio request body: %v", err)
		http.Error(rw, "Error while decoding portfolio data", http.StatusBadRequest)
		return
	}

	handler.PortfolioStore.Add(context.Background(), &portfolio)
}
