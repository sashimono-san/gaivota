package data

type Order struct {
	ID         int
	SKU        string
	PositionId int
	Amount     float32
	UnitPrice  float32
	TotalPrice float32
	Type       string // Must be `sell` or `buy`
	Exchange   string
	CreatedOn  string
	UpdatedOn  string
	DeletedOn  string
}

var orders = []*Order{}
