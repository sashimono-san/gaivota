package gaivota

import "context"

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type UserStore interface {
	// Add creates a new User in the UsersStore and returns User with ID
	Add(context.Context, *User) (*User, error)
	// Returns all users in the store
	All(context.Context) (*[]User, error)
	// Delete the User from the store
	Delete(ctx context.Context, id int) error
	// Gets User if `ID` exists
	Get(ctx context.Context, id int) (*User, error)
	// Update the User in the store.
	Update(context.Context, *User) error
}

type Portfolio struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user"`
	Name      string `json:"name"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type PortfolioStore interface {
	// Add creates a new Portfolio in the PortfoliosStore and returns Portfolio with ID
	Add(context.Context, *Portfolio) (*Portfolio, error)
	// Returns all Portfolios in the store
	All(context.Context) (*[]Portfolio, error)
	// Delete the Portfolio from the store
	Delete(ctx context.Context, id int) error
	// Gets Portfolio if `ID` exists
	Get(ctx context.Context, id int) (*Portfolio, error)
	// Gets all Portfolios for user
	GetByUserID(ctx context.Context, userId int) (*[]Portfolio, error)
	// Update the Portfolio in the store.
	Update(context.Context, *Portfolio) error
}

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

type WalletStore interface {
	// Add creates a new Wallet in the WalletsStore and returns Wallet with ID
	Add(context.Context, *Wallet) (*Wallet, error)
	// Returns all Wallets in the store
	All(context.Context) (*[]Wallet, error)
	// Delete the Wallet from the store
	Delete(ctx context.Context, id int) error
	// Gets Wallet if `ID` exists
	Get(ctx context.Context, id int) (*Wallet, error)
	// Gets all Wallets for user
	GetByUserID(ctx context.Context, userId int) (*[]Wallet, error)
	// Update the Wallet in the store.
	Update(context.Context, *Wallet) error
}

type Investment struct {
	ID          int    `json:"id"`
	PortfolioID int    `json:"portfolio"`
	Token       string `json:"token"`
	TokenSymbol string `json:"symbol"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

type InvestmentStore interface {
	// Add creates a new Investment in the InvestmentsStore and returns Investment with ID
	Add(context.Context, *Investment) (*Investment, error)
	// Returns all Investments in the store
	All(context.Context) (*[]Investment, error)
	// Delete the Investment from the store
	Delete(ctx context.Context, id int) error
	// Gets Investment if `ID` exists
	Get(ctx context.Context, id int) (*Investment, error)
	// Gets all Investments for user
	GetByUserID(ctx context.Context, userId int) (*[]Investment, error)
	// Gets all Investments for portfolio
	GetByPortfolioID(ctx context.Context, portfolioId int) (*[]Investment, error)
	// Update the Investment in the store.
	Update(context.Context, *Investment) error
}

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

type PositionStore interface {
	// Add creates a new Position in the PositionsStore and returns Position with ID
	Add(context.Context, *Position) (*Position, error)
	// Returns all Positions in the store
	All(context.Context) ([]Position, error)
	// Delete the Position from the store
	Delete(ctx context.Context, id int) error
	// Gets Position if `ID` exists
	Get(ctx context.Context, id int) (*Position, error)
	// Update the Position in the store.
	Update(context.Context, *Position) error
}

type Holding struct {
	ID         int      `json:"id"`
	WalletID   int      `json:"wallet"`
	Wallet     Wallet   `json:"-"`
	PositionID int      `json:"position"`
	Position   Position `json:"-"`
	Amount     float64  `json:"amount"`
	CreatedAt  string   `json:"-"`
	UpdatedAt  string   `json:"-"`
	DeletedAt  string   `json:"-"`
}

type HoldingStore interface {
	// Add creates a new Holding in the HoldingsStore and returns Holding with ID
	Add(context.Context, *Holding) (*Holding, error)
	// Returns all Holdings in the store
	All(context.Context) ([]Holding, error)
	// Delete the Holding from the store
	Delete(ctx context.Context, id int) error
	// Gets Holding if `ID` exists
	Get(ctx context.Context, id int) (*Holding, error)
	// Gets all Holdings for user
	GetByUserID(ctx context.Context, userId int) (*[]Holding, error)
	// Gets all Holdings for wallet
	GetByWalletID(ctx context.Context, walletId int) (*[]Holding, error)
	// Gets all Holdings for position
	GetByPositionID(ctx context.Context, positionId int) (*[]Holding, error)
	// Update the Holding in the store.
	Update(context.Context, *Holding) error
}

// Operations enum
type OrderOperation string

const (
	OrderOperationSell OrderOperation = "sell"
	OrderOperationBuy  OrderOperation = "buy"
)

// Order type enum
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
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

type OrderStore interface {
	// Add creates a new Order in the OrdersStore and returns Order with ID
	Add(context.Context, *Order) (*Order, error)
	// Returns all Orders in the store
	All(context.Context) ([]Order, error)
	// Delete the Order from the store
	Delete(ctx context.Context, id int) error
	// Gets Order if `ID` exists
	Get(ctx context.Context, id int) (*Order, error)
	// Update the Order in the store.
	Update(context.Context, *Order) error
}
