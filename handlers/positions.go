package handlers

import (
	"gaivota/data"
	"log"
	"net/http"
)

type Positions struct {
	logger *log.Logger
}

func NewPosition(logger *log.Logger) *Positions {
	return &Positions{logger}
}

func (handler *Positions) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		handler.GetPositions(res, req)
		return
	}

	if req.Method == http.MethodPost {
		handler.AddPosition(res, req)
		return
	}

	// if req.Method == http.MethodPut {
	// 	req.URL.Path
	// 	handler.UpdatePosition(res, req)
	// 	return
	// }

	res.WriteHeader(http.StatusMethodNotAllowed)
}

func (handler *Positions) GetPositions(res http.ResponseWriter, req *http.Request) {
	handler.logger.Println("Handle GET Positions")

	positions := data.GetPositions()
	err := positions.ToJSON(res)

	if err != nil {
		http.Error(res, "Unable to parse positions' list", http.StatusInternalServerError)
	}
}

func (handler *Positions) AddPosition(res http.ResponseWriter, req *http.Request) {
	handler.logger.Println("Handle POST Positions")

	position := &data.Position{}
	err := position.FromJSON(req.Body)

	if err != nil {
		http.Error(res, "Unable to parse position", http.StatusBadRequest)
	}

	// handler.logger.Printf("Position: %#v", position)
	data.AddPosition(position)
}
