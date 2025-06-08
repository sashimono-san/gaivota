package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/leoschet/gaivota"
	"github.com/leoschet/gaivota/internal/config"
	"github.com/leoschet/gaivota/log"
	"github.com/leoschet/gaivota/postgres"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	logger := log.New("Gaivota-CLI - ")

	rootPath, err := os.Getwd()
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, "Error getting root path: %v", err)
	}

	// Try to read local config first, fallback to regular config
	configPath := path.Join(rootPath, "/config.local.json")
	settings, err := config.ReadFile(configPath)
	if err != nil {
		// Fallback to regular config
		configPath = path.Join(rootPath, "/config.json")
		settings, err = config.ReadFile(configPath)
		if err != nil {
			logger.Log(gaivota.LogLevelFatal, "Error while reading config file: %v", err)
		}
		
		// If using regular config, check if we need to adjust the database host for local usage
		if settings.DatabaseConnString != "" {
			// Replace @db:5432 with @localhost:5555 for local development
			settings.DatabaseConnString = strings.Replace(settings.DatabaseConnString, "@db:5432", "@localhost:5555", 1)
			logger.Log(gaivota.LogLevelInfo, "Adjusted database connection for local development")
		}
	}

	db, err := postgres.Connect(context.Background(), settings.DatabaseConnString)
	if err != nil {
		logger.Log(gaivota.LogLevelFatal, "Error while connecting to Postgres: %v", err)
	}
	defer db.Close()

	pgClient := db.NewPostgresClient()

	command := os.Args[1]
	switch command {
	case "users":
		handleUsers(pgClient, os.Args[2:])
	case "portfolios":
		handlePortfolios(pgClient, os.Args[2:])
	case "wallets":
		handleWallets(pgClient, os.Args[2:])
	case "investments":
		handleInvestments(pgClient, os.Args[2:])
	case "positions":
		handlePositions(pgClient, os.Args[2:])
	case "orders":
		handleOrders(pgClient, os.Args[2:])
	case "health":
		handleHealth(db)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Gaivota CLI - Portfolio Management Tool")
	fmt.Println("")
	fmt.Println("Usage: gaivota-cli <command> [args]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  health                    Check database connection")
	fmt.Println("  users <subcommand>        Manage users")
	fmt.Println("    list                    List all users")
	fmt.Println("    get <id>                Get user by ID")
	fmt.Println("    create <email> <first> <last>  Create new user")
	fmt.Println("  portfolios <subcommand>   Manage portfolios")
	fmt.Println("    list                    List all portfolios")
	fmt.Println("    list-by-user <user_id>  List portfolios for user")
	fmt.Println("    get <id>                Get portfolio by ID")
	fmt.Println("    create <user_id> <name> Create new portfolio")
	fmt.Println("  wallets <subcommand>      Manage wallets")
	fmt.Println("    list                    List all wallets")
	fmt.Println("    list-by-user <user_id>  List wallets for user")
	fmt.Println("    get <id>                Get wallet by ID")
	fmt.Println("  investments <subcommand>  Manage investments")
	fmt.Println("    list                    List all investments")
	fmt.Println("    get <id>                Get investment by ID")
	fmt.Println("  positions <subcommand>    Manage positions")
	fmt.Println("    list                    List all positions")
	fmt.Println("    get <id>                Get position by ID")
	fmt.Println("  orders <subcommand>       Manage orders")
	fmt.Println("    list                    List all orders")
	fmt.Println("    get <id>                Get order by ID")
}

func handleHealth(db gaivota.HealthChecker) {
	msg, err := db.Ping()
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Database connection healthy: %s\n", msg)
}

func handleUsers(client *gaivota.Client, args []string) {
	ctx := context.Background()
	
	if len(args) == 0 {
		fmt.Println("Missing subcommand for users")
		return
	}

	switch args[0] {
	case "list":
		users, err := client.UserStore.All(ctx)
		if err != nil {
			fmt.Printf("Error listing users: %v\n", err)
			return
		}
		
		fmt.Println("Users:")
		fmt.Printf("%-5s %-25s %-15s %-15s\n", "ID", "Email", "First Name", "Last Name")
		fmt.Println("-------------------------------------------------------------")
		for _, user := range *users {
			fmt.Printf("%-5d %-25s %-15s %-15s\n", user.ID, user.Email, user.FirstName, user.LastName)
		}

	case "get":
		if len(args) < 2 {
			fmt.Println("Missing user ID")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid user ID: %s\n", args[1])
			return
		}
		
		user, err := client.UserStore.Get(ctx, id)
		if err != nil {
			fmt.Printf("Error getting user: %v\n", err)
			return
		}
		
		fmt.Printf("User Details:\n")
		fmt.Printf("  ID: %d\n", user.ID)
		fmt.Printf("  Email: %s\n", user.Email)
		fmt.Printf("  Name: %s %s\n", user.FirstName, user.LastName)
		fmt.Printf("  Created: %s\n", user.CreatedAt)

	case "create":
		if len(args) < 4 {
			fmt.Println("Usage: users create <email> <first_name> <last_name>")
			return
		}
		
		user := &gaivota.User{
			Email:     args[1],
			FirstName: args[2],
			LastName:  args[3],
		}
		
		createdUser, err := client.UserStore.Add(ctx, user)
		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			return
		}
		
		fmt.Printf("User created successfully:\n")
		fmt.Printf("  ID: %d\n", createdUser.ID)
		fmt.Printf("  Email: %s\n", createdUser.Email)
		fmt.Printf("  Name: %s %s\n", createdUser.FirstName, createdUser.LastName)

	default:
		fmt.Printf("Unknown users subcommand: %s\n", args[0])
	}
}

func handlePortfolios(client *gaivota.Client, args []string) {
	ctx := context.Background()
	
	if len(args) == 0 {
		fmt.Println("Missing subcommand for portfolios")
		return
	}

	switch args[0] {
	case "list":
		portfolios, err := client.PortfolioStore.All(ctx)
		if err != nil {
			fmt.Printf("Error listing portfolios: %v\n", err)
			return
		}
		
		fmt.Println("Portfolios:")
		fmt.Printf("%-5s %-10s %-30s\n", "ID", "User ID", "Name")
		fmt.Println("-----------------------------------------------")
		for _, portfolio := range *portfolios {
			fmt.Printf("%-5d %-10d %-30s\n", portfolio.ID, portfolio.UserID, portfolio.Name)
		}

	case "list-by-user":
		if len(args) < 2 {
			fmt.Println("Missing user ID")
			return
		}
		userID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid user ID: %s\n", args[1])
			return
		}
		
		portfolios, err := client.PortfolioStore.GetByUserID(ctx, userID)
		if err != nil {
			fmt.Printf("Error listing portfolios for user: %v\n", err)
			return
		}
		
		fmt.Printf("Portfolios for User %d:\n", userID)
		fmt.Printf("%-5s %-30s\n", "ID", "Name")
		fmt.Println("------------------------------------")
		for _, portfolio := range *portfolios {
			fmt.Printf("%-5d %-30s\n", portfolio.ID, portfolio.Name)
		}

	case "get":
		if len(args) < 2 {
			fmt.Println("Missing portfolio ID")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid portfolio ID: %s\n", args[1])
			return
		}
		
		portfolio, err := client.PortfolioStore.Get(ctx, id)
		if err != nil {
			fmt.Printf("Error getting portfolio: %v\n", err)
			return
		}
		
		fmt.Printf("Portfolio Details:\n")
		fmt.Printf("  ID: %d\n", portfolio.ID)
		fmt.Printf("  User ID: %d\n", portfolio.UserID)
		fmt.Printf("  Name: %s\n", portfolio.Name)
		fmt.Printf("  Created: %s\n", portfolio.CreatedAt)

	case "create":
		if len(args) < 3 {
			fmt.Println("Usage: portfolios create <user_id> <name>")
			return
		}
		
		userID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid user ID: %s\n", args[1])
			return
		}
		
		portfolio := &gaivota.Portfolio{
			UserID: userID,
			Name:   args[2],
		}
		
		createdPortfolio, err := client.PortfolioStore.Add(ctx, portfolio)
		if err != nil {
			fmt.Printf("Error creating portfolio: %v\n", err)
			return
		}
		
		fmt.Printf("Portfolio created successfully:\n")
		fmt.Printf("  ID: %d\n", createdPortfolio.ID)
		fmt.Printf("  User ID: %d\n", createdPortfolio.UserID)
		fmt.Printf("  Name: %s\n", createdPortfolio.Name)

	default:
		fmt.Printf("Unknown portfolios subcommand: %s\n", args[0])
	}
}

func handleWallets(client *gaivota.Client, args []string) {
	ctx := context.Background()
	
	if len(args) == 0 {
		fmt.Println("Missing subcommand for wallets")
		return
	}

	switch args[0] {
	case "list":
		wallets, err := client.WalletStore.All(ctx)
		if err != nil {
			fmt.Printf("Error listing wallets: %v\n", err)
			return
		}
		
		fmt.Println("Wallets:")
		fmt.Printf("%-5s %-10s %-20s %-15s %-40s\n", "ID", "User ID", "Name", "Total Value", "Address")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, wallet := range *wallets {
			fmt.Printf("%-5d %-10d %-20s $%-14.2f %-40s\n", wallet.ID, wallet.UserID, wallet.Name, wallet.TotalValue, wallet.Address)
		}

	case "list-by-user":
		if len(args) < 2 {
			fmt.Println("Missing user ID")
			return
		}
		userID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid user ID: %s\n", args[1])
			return
		}
		
		wallets, err := client.WalletStore.GetByUserID(ctx, userID)
		if err != nil {
			fmt.Printf("Error listing wallets for user: %v\n", err)
			return
		}
		
		fmt.Printf("Wallets for User %d:\n", userID)
		fmt.Printf("%-5s %-20s %-15s %-40s\n", "ID", "Name", "Total Value", "Address")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, wallet := range *wallets {
			fmt.Printf("%-5d %-20s $%-14.2f %-40s\n", wallet.ID, wallet.Name, wallet.TotalValue, wallet.Address)
		}

	case "get":
		if len(args) < 2 {
			fmt.Println("Missing wallet ID")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid wallet ID: %s\n", args[1])
			return
		}
		
		wallet, err := client.WalletStore.Get(ctx, id)
		if err != nil {
			fmt.Printf("Error getting wallet: %v\n", err)
			return
		}
		
		fmt.Printf("Wallet Details:\n")
		fmt.Printf("  ID: %d\n", wallet.ID)
		fmt.Printf("  User ID: %d\n", wallet.UserID)
		fmt.Printf("  Name: %s\n", wallet.Name)
		fmt.Printf("  Total Value: $%.2f\n", wallet.TotalValue)
		fmt.Printf("  Address: %s\n", wallet.Address)
		fmt.Printf("  Location: %s\n", wallet.Location)
		fmt.Printf("  Created: %s\n", wallet.CreatedAt)

	default:
		fmt.Printf("Unknown wallets subcommand: %s\n", args[0])
	}
}

func handleInvestments(client *gaivota.Client, args []string) {
	ctx := context.Background()
	
	if len(args) == 0 {
		fmt.Println("Missing subcommand for investments")
		return
	}

	switch args[0] {
	case "list":
		investments, err := client.InvestmentStore.All(ctx)
		if err != nil {
			fmt.Printf("Error listing investments: %v\n", err)
			return
		}
		
		fmt.Println("Investments:")
		fmt.Printf("%-5s %-15s %-20s %-10s\n", "ID", "Portfolio ID", "Token", "Symbol")
		fmt.Println("----------------------------------------------------")
		for _, investment := range *investments {
			fmt.Printf("%-5d %-15d %-20s %-10s\n", investment.ID, investment.PortfolioID, investment.Token, investment.TokenSymbol)
		}

	case "get":
		if len(args) < 2 {
			fmt.Println("Missing investment ID")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid investment ID: %s\n", args[1])
			return
		}
		
		investment, err := client.InvestmentStore.Get(ctx, id)
		if err != nil {
			fmt.Printf("Error getting investment: %v\n", err)
			return
		}
		
		fmt.Printf("Investment Details:\n")
		fmt.Printf("  ID: %d\n", investment.ID)
		fmt.Printf("  Portfolio ID: %d\n", investment.PortfolioID)
		fmt.Printf("  Token: %s\n", investment.Token)
		fmt.Printf("  Symbol: %s\n", investment.TokenSymbol)
		fmt.Printf("  Created: %s\n", investment.CreatedAt)

	default:
		fmt.Printf("Unknown investments subcommand: %s\n", args[0])
	}
}

func handlePositions(client *gaivota.Client, args []string) {
	ctx := context.Background()
	
	if len(args) == 0 {
		fmt.Println("Missing subcommand for positions")
		return
	}

	switch args[0] {
	case "list":
		positions, err := client.PositionStore.All(ctx)
		if err != nil {
			fmt.Printf("Error listing positions: %v\n", err)
			return
		}
		
		fmt.Println("Positions:")
		fmt.Printf("%-5s %-15s %-15s %-15s %-15s\n", "ID", "Investment ID", "Amount", "Avg Price", "Profit")
		fmt.Println("-----------------------------------------------------------------------")
		for _, position := range *positions {
			fmt.Printf("%-5d %-15d %-15.6f $%-14.2f $%-14.2f\n", 
				position.ID, position.InvestmentID, position.Amount, position.AveragePrice, position.Profit)
		}

	case "get":
		if len(args) < 2 {
			fmt.Println("Missing position ID")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid position ID: %s\n", args[1])
			return
		}
		
		position, err := client.PositionStore.Get(ctx, id)
		if err != nil {
			fmt.Printf("Error getting position: %v\n", err)
			return
		}
		
		fmt.Printf("Position Details:\n")
		fmt.Printf("  ID: %d\n", position.ID)
		fmt.Printf("  Investment ID: %d\n", position.InvestmentID)
		fmt.Printf("  Amount: %.6f\n", position.Amount)
		fmt.Printf("  Average Price: $%.2f\n", position.AveragePrice)
		fmt.Printf("  Profit: $%.2f\n", position.Profit)
		fmt.Printf("  Created: %s\n", position.CreatedAt)

	default:
		fmt.Printf("Unknown positions subcommand: %s\n", args[0])
	}
}

func handleOrders(client *gaivota.Client, args []string) {
	ctx := context.Background()
	
	if len(args) == 0 {
		fmt.Println("Missing subcommand for orders")
		return
	}

	switch args[0] {
	case "list":
		orders, err := client.OrderStore.All(ctx)
		if err != nil {
			fmt.Printf("Error listing orders: %v\n", err)
			return
		}
		
		fmt.Println("Orders:")
		fmt.Printf("%-5s %-12s %-10s %-12s %-12s %-8s %-8s %-15s\n", 
			"ID", "Position ID", "Amount", "Unit Price", "Total", "Op", "Type", "Exchange")
		fmt.Println("-----------------------------------------------------------------------------------")
		for _, order := range orders {
			fmt.Printf("%-5d %-12d %-10.4f $%-11.2f $%-11.2f %-8s %-8s %-15s\n", 
				order.ID, order.PositionID, order.Amount, order.UnitPrice, order.TotalPrice,
				order.Operation, order.Type, order.Exchange)
		}

	case "get":
		if len(args) < 2 {
			fmt.Println("Missing order ID")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Invalid order ID: %s\n", args[1])
			return
		}
		
		order, err := client.OrderStore.Get(ctx, id)
		if err != nil {
			fmt.Printf("Error getting order: %v\n", err)
			return
		}
		
		fmt.Printf("Order Details:\n")
		fmt.Printf("  ID: %d\n", order.ID)
		fmt.Printf("  Position ID: %d\n", order.PositionID)
		fmt.Printf("  Amount: %.4f\n", order.Amount)
		fmt.Printf("  Unit Price: $%.2f\n", order.UnitPrice)
		fmt.Printf("  Total Price: $%.2f\n", order.TotalPrice)
		fmt.Printf("  Operation: %s\n", order.Operation)
		fmt.Printf("  Type: %s\n", order.Type)
		fmt.Printf("  Exchange: %s\n", order.Exchange)
		fmt.Printf("  Executed At: %s\n", order.ExecutedAt)
		fmt.Printf("  Created: %s\n", order.CreatedAt)

	default:
		fmt.Printf("Unknown orders subcommand: %s\n", args[0])
	}
}
