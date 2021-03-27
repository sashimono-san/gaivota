package data

type Investment struct {
	ID          int    `json:"id"`
	PortfolioID int    `json:"portfolio"`
	Token       string `json:"token"`
	TokenSymbol string `json:"symbol"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

var investments = []*Investment{}
