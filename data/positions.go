package data

import (
	"encoding/json"
	"io"
	"time"
)

type Position struct {
	ID           int     `json:"id"`
	InvestmentID int     `json:"investment"`
	Amount       float64 `json:"amount"`
	AveragePrice float64 `json:"averagePrice"`
	Profit       float64 `json:"profit,omitempty"`
	CreatedAt    string  `json:"-"`
	UpdatedAt    string  `json:"-"`
	DeletedAt    string  `json:"-"`
}

func (position *Position) FromJSON(r io.Reader) error {
	// https://golang.org/pkg/encoding/json/
	decoder := json.NewDecoder(r)
	return decoder.Decode(position)
}

type Positions []*Position

func (positions *Positions) ToJSON(w io.Writer) error {
	// https://golang.org/pkg/encoding/json/
	encoder := json.NewEncoder(w)
	return encoder.Encode(positions)
}

func GetPositions() Positions {
	return positions
}

func AddPosition(position *Position) {
	positions = append(positions, position)
}

var positions = Positions{
	&Position{
		ID:           1,
		InvestmentID: 1,
		Amount:       1.00,
		AveragePrice: 1.681,
		CreatedAt:    time.Now().UTC().String(),
		UpdatedAt:    time.Now().UTC().String(),
	},
}
