package data

import (
	"encoding/json"
	"io"
	"time"
)

type Position struct {
	ID           int     `json:"-"`
	UUID         string  `json:"uuid"`
	InvestmentId int     `json:"investmentUuid"`
	Amount       float32 `json:"amount"`
	AveragePrice float32 `json:"averagePrice"`
	Profit       float32 `json:"profit,omitempty"`
	CreatedOn    string  `json:"-"`
	UpdatedOn    string  `json:"-"`
	DeletedOn    string  `json:"-"`
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
		UUID:         "abc123",
		InvestmentId: 1,
		Amount:       1.00,
		AveragePrice: 1.681,
		CreatedOn:    time.Now().UTC().String(),
		UpdatedOn:    time.Now().UTC().String(),
	},
}
