package data

import (
	"database/sql/driver"
	"errors"
)

// Operations enum
type OrderOperation string

func (operation *OrderOperation) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Operation source is not []byte")
	}
	*operation = OrderOperation(string(asBytes))
	return nil
}

func (operation OrderOperation) Value() (driver.Value, error) {
	return string(operation), nil
}

const (
	sell OrderOperation = "sell"
	buy  OrderOperation = "buy"
)

// Order type enum
type OrderType string

func (orderType *OrderType) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Order type source is not []byte")
	}
	*orderType = OrderType(string(asBytes))
	return nil
}

func (orderType OrderType) Value() (driver.Value, error) {
	return string(orderType), nil
}

const (
	limit  OrderType = "limit"
	market OrderType = "market"
)

type Order struct {
	ID         int            `json:"id"`
	PositionID int            `json:"position"`
	Amount     float32        `json:"amount"`
	UnitPrice  float32        `json:"unitPrice"`
	TotalPrice float32        `json:"totalPrice"`
	Operation  OrderOperation `json:"operation"`
	Type       OrderType      `json:"type"`
	Exchange   string         `json:"exchange"`
	ExecutedAt string         `json:"executedAt"`
	CreatedAt  string         `json:"-"`
	UpdatedAt  string         `json:"-"`
	DeletedAt  string         `json:"-"`
}

var orders = []*Order{}
