package gaivota

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
	Add(*User) (*User, error)
	// Returns all users in the store
	All() ([]User, error)
	// Delete the User from the store
	Delete(*User) error
	// Get retrieves User if `ID` exists
	Get(id int) (*User, error)
	// Update the User in the store.
	Update(*User)
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
	Add(*Portfolio) (*Portfolio, error)
	// Returns all Portfolios in the store
	All() ([]Portfolio, error)
	// Delete the Portfolio from the store
	Delete(*Portfolio) error
	// Get retrieves Portfolio if `ID` exists
	Get(id int) (*Portfolio, error)
	// Update the Portfolio in the store.
	Update(*Portfolio)
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
	Add(*Wallet) (*Wallet, error)
	// Returns all Wallets in the store
	All() ([]Wallet, error)
	// Delete the Wallet from the store
	Delete(*Wallet) error
	// Get retrieves Wallet if `ID` exists
	Get(id int) (*Wallet, error)
	// Update the Wallet in the store.
	Update(*Wallet)
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
	Add(*Investment) (*Investment, error)
	// Returns all Investments in the store
	All() ([]Investment, error)
	// Delete the Investment from the store
	Delete(*Investment) error
	// Get retrieves Investment if `ID` exists
	Get(id int) (*Investment, error)
	// Update the Investment in the store.
	Update(*Investment)
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
	Add(*Position) (*Position, error)
	// Returns all Positions in the store
	All() ([]Position, error)
	// Delete the Position from the store
	Delete(*Position) error
	// Get retrieves Position if `ID` exists
	Get(id int) (*Position, error)
	// Update the Position in the store.
	Update(*Position)
}

type Holding struct {
	ID         int     `json:"id"`
	WalletID   int     `json:"wallet"`
	PositionID int     `json:"position"`
	Amount     float64 `json:"amount"`
	CreatedAt  string  `json:"-"`
	UpdatedAt  string  `json:"-"`
	DeletedAt  string  `json:"-"`
}

type HoldingStore interface {
	// Add creates a new Holding in the HoldingsStore and returns Holding with ID
	Add(*Holding) (*Holding, error)
	// Returns all Holdings in the store
	All() ([]Holding, error)
	// Delete the Holding from the store
	Delete(*Holding) error
	// Get retrieves Holding if `ID` exists
	Get(id int) (*Holding, error)
	// Update the Holding in the store.
	Update(*Holding)
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
	Add(*Order) (*Order, error)
	// Returns all Orders in the store
	All() ([]Order, error)
	// Delete the Order from the store
	Delete(*Order) error
	// Get retrieves Order if `ID` exists
	Get(id int) (*Order, error)
	// Update the Order in the store.
	Update(*Order)
}
