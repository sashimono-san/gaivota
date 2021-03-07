package data

type Investment struct {
	ID          int
	SKU         string
	PortfolioId int
	Token       string
	TokenSymbol string
	// Amount       float32
	// AveragePrice float32
	// Profit       float32
	CreatedOn string
	UpdatedOn string
	DeletedOn string
}

var investments = []*Investment{}
