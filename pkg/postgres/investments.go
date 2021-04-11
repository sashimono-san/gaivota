package postgres

import (
	"context"
	"fmt"
	"gaivota"

	"github.com/jackc/pgx/v4/pgx"
)

type InvestmentStore struct {
	Database *Database
}

func (store *InvestmentStore) getByFK(ctx context.Context, fk_column string, fk int) (*[]gaivota.Investment, error) {
	query := `select "id", "portfolio_id", "token", "token_symbol", "created_at", "updated_at", "deleted_at"
						from investments where $1 = $2`

	rows, err := store.Database.Pool.Query(ctx, query, fk_column, fk)

	if err != nil {
		return nil, fmt.Errorf("Could not get holdings where %s is %v: %w", fk_column, fk, err)
	}

	return store.scanAll(rows)
}

func (store *InvestmentStore) scanAll(rows pgx.Rows) (*[]gaivota.Investment, error) {
	var investments []gaivota.Investment

	for rows.Next() {
		investment, err := store.scanOne(rows)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning investments: %w", err)
		}

		investments = append(investments, *investment)
	}

	return &investments, nil
}

func (store *InvestmentStore) scanOne(row pgx.Row) (*gaivota.Investment, error) {
	var investment gaivota.Investment

	err := row.Scan(
		&investment.ID, &investment.PortfolioID, &investment.Token, &investment.TokenSymbol,
		&investment.CreatedAt, &investment.UpdatedAt, &investment.DeletedAt,
	)

	return &investment, err
}

func (store *InvestmentStore) Add(ctx context.Context, investment *gaivota.Investment) (*gaivota.Investment, error) {
	query := `insert into investments ("portfolio_id", "token", "token_symbol")
						values ($1, $2)
						returning "id", "portfolio_id", "token", "token_symbol", "created_at", "updated_at", "deleted_at"`

	row := store.Database.Pool.QueryRow(ctx, query, investment.PortfolioID, investment.Token, investment.TokenSymbol)

	newInvestment, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not insert investment of %s in portfolio %v: %w", investment.Token, investment.PortfolioID, err)
	}

	return newInvestment, nil
}

func (store *InvestmentStore) All(ctx context.Context) (*[]gaivota.Investment, error) {
	query := `select "id", "portfolio_id", "token", "token_symbol", "created_at", "updated_at", "deleted_at"
						from investments`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get investments: %w", err)
	}

	return store.scanAll(rows)
}

func (store *InvestmentStore) Delete(ctx context.Context, id int) error {
	query := `update investments
						set deleted_at = now(),
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, id,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete investment %v: %w", id, err)
	}

	return nil
}

func (store *InvestmentStore) Get(ctx context.Context, id int) (*gaivota.Investment, error) {
	query := `select "id", "portfolio_id", "token", "token_symbol", "created_at", "updated_at", "deleted_at"
						from investments where id = $1`

	row := store.Database.Pool.QueryRow(ctx, query, id)

	investment, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not get investment %v: %w", id, err)
	}

	return investment, nil
}

func (store *InvestmentStore) GetByUserID(ctx context.Context, userId int) (*[]gaivota.Investment, error) {
	return store.getByFK(ctx, "user_id", userId)
}

func (store *InvestmentStore) GetByPortfolioID(ctx context.Context, portfolioId int) (*[]gaivota.Investment, error) {
	return store.getByFK(ctx, "portfolio_id", portfolioId)
}

func (store *InvestmentStore) Update(ctx context.Context, investment *gaivota.Investment) error {
	query := `update investments
						set portfolio_id = $1
						where id = $2`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, &investment.PortfolioID, &investment.ID,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update investment %v: %w", investment.ID, err)
	}

	return nil
}
