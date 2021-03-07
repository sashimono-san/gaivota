package data

type Wallet struct {
	ID         int
	SKU        string
	UserId     int
	PositionId int
	Name       string
	Amount     float32
	Address    string
	Location   string
	CreatedOn  string
	UpdatedOn  string
	DeletedOn  string
}

var wallets = []*Wallet{}
