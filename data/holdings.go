package data

type Holding struct {
	ID         int     `json:"id"`
	PositionID int     `json:"position"`
	WalletID   int     `json:"wallet"` // TODO: Wallet has positions or position has wallet? maybe many to many?
	Amount     float64 `json:"amount"`
	CreatedAt  string  `json:"-"`
	UpdatedAt  string  `json:"-"`
	DeletedAt  string  `json:"-"`
}
