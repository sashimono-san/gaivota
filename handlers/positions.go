package handlers

import (
	"gaivota/data"
	"log"
	"net/http"
)

type Positions struct {
	logger *log.Logger
}

func NewPositions(logger *log.Logger) *Positions {
	return &Positions{logger}
}

func (handler *Positions) Get(res http.ResponseWriter, req *http.Request) {
	handler.logger.Println("Handle GET Positions")

	positions := data.GetPositions()
	err := positions.ToJSON(res)

	if err != nil {
		http.Error(res, "Unable to parse positions' list", http.StatusInternalServerError)
	}
}

func (handler *Positions) Add(res http.ResponseWriter, req *http.Request) {
	handler.logger.Println("Handle POST Positions")

	position := &data.Position{}
	err := position.FromJSON(req.Body)

	if err != nil {
		http.Error(res, "Unable to parse position", http.StatusBadRequest)
	}

	// handler.logger.Printf("Position: %#v", position)
	data.AddPosition(position)
}
