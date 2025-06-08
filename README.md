# Gaivota 🐦

**Portfolio Management Tool**

Gaivota (Portuguese for "seagull") is a comprehensive portfolio management API written in Go. It's designed to help users track and manage their investment portfolios across multiple wallets and exchanges.

## Features

- **Multi-Wallet Support**: Track assets across different wallets and exchanges
- **Investment Tracking**: Monitor positions with average prices and automatic profit/loss calculations
- **Order Management**: Record buy/sell orders with support for limit and market orders
- **Multi-User**: Support for multiple users with individual portfolios
- **Real-time Health Checks**: Built-in health monitoring endpoints
- **Soft Deletes**: Data integrity with soft delete functionality

## Architecture

### Tech Stack
- **Backend**: Go 1.16 with custom HTTP router
- **Database**: PostgreSQL 13 with migrations
- **Containerization**: Docker with docker-compose
- **Database Driver**: pgx/v4 for PostgreSQL connectivity

### Domain Model

The system follows a well-structured domain model with these core entities:

1. **Users** - Account holders of the system
2. **Portfolios** - Logical groupings of investments per user
3. **Wallets** - Physical or digital storage locations with addresses
4. **Investments** - Specific tokens/assets within portfolios
5. **Positions** - Investment amounts with average prices and profit tracking
6. **Holdings** - Relationships between positions and wallets (where assets are stored)
7. **Orders** - Buy/sell transactions with exchange information

### Project Structure

```
├── cmd/gaivota/          # Application entry point
├── handlers/             # HTTP request handlers
├── internal/config/      # Configuration management
├── log/                  # Custom logging
├── mux/                  # HTTP routing and endpoints
├── postgres/             # Database layer implementations
├── migrations/           # Database schema migrations
└── gaivota.go           # Core domain types and interfaces
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.16+ (for local development)

### Running with Docker

1. **Create configuration file**
   ```bash
   cp config.exemple.json config.json
   ```

2. **Start the services**
   ```bash
   docker compose up -d
   ```

   This will start:
   - API server on port 8888
   - PostgreSQL database on port 5555

### Configuration

The application uses a JSON configuration file with the following structure:

```json
{
  "Port": 9090,
  "DatabaseConnString": "postgres://gaivota:secretpassword@db:5432/gaivota"
}
```

### Database

The application uses PostgreSQL with automated migrations. The database schema includes:

- **users**: User account information
- **portfolios**: Investment portfolio groupings
- **wallets**: Asset storage locations
- **investments**: Tracked tokens/assets
- **positions**: Investment amounts and pricing
- **holdings**: Position-wallet relationships
- **orders**: Transaction history

All tables include automatic timestamp tracking and soft delete functionality.

### Available Interfaces

**1. REST API Server**
- Health checks (`/health`)
- Portfolio management endpoints
- User management endpoints
- Investment tracking endpoints
- Order processing endpoints

**2. Command Line Interface (CLI)**
- Direct database access for all entities
- User-friendly commands for data management
- Perfect for administration and testing

### Development

#### Building Applications

The project follows Go's standard `cmd/` pattern with multiple entry points:

```bash
# Build API server
go build -o bin/gaivota ./cmd/gaivota

# Build CLI tool
go build -o bin/gaivota-cli ./cmd/gaivota-cli

# Or build both
go build -o bin/ ./cmd/...
```

#### Running Locally

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Setup configuration**
   ```bash
   cp config.example.json config.json
   ```

3. **Start database**
   ```bash
   docker compose up -d db
   ```

4. **Run database migrations**
   ```bash
   # Configure tern.conf based on tern.example.conf
   tern migrate
   ```

5. **Start the API server**
   ```bash
   go run cmd/gaivota/main.go
   # or
   ./bin/gaivota
   ```

6. **Use the CLI tool**
   ```bash
   go run cmd/gaivota-cli/main.go --help
   # or
   ./bin/gaivota-cli --help
   ```

### CLI Usage Examples

```bash
# Check database connection
./gaivota-cli health

# List all users
./gaivota-cli users list

# Create a new user
./gaivota-cli users create "john@example.com" "John" "Doe"

# List portfolios for a user
./gaivota-cli portfolios list-by-user 1

# Create a portfolio
./gaivota-cli portfolios create 1 "My Crypto Portfolio"

# List all wallets
./gaivota-cli wallets list

# Get specific investment details
./gaivota-cli investments get 1
```

## Database Schema

The system uses PostgreSQL with the following key relationships:

- Users → Portfolios (1:many)
- Users → Wallets (1:many)
- Portfolios → Investments (1:many)
- Investments → Positions (1:many)
- Positions → Orders (1:many)
- Positions ↔ Wallets (many:many through Holdings)

## Contributing

The project follows Go best practices with:
- Interface-based design for testability
- Clean separation of concerns
- Dependency injection
- Graceful shutdown handling
- Comprehensive error handling and logging

## License

See LICENSE file for details.
