package data

type Wallet struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user"`
	Name       string  `json:"name"`
	TotalValue float32 `json:"totalValue"`
	Address    string  `json:"address"`
	Location   string  `json:"location"`
	CreatedAt  string  `json:"-"`
	UpdatedAt  string  `json:"-"`
	DeletedAt  string  `json:"-"`
}

var wallets = []*Wallet{}
