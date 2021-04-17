package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/leoschet/gaivota"
)

type PortfolioStore struct {
	Database *Database
}

func (store *PortfolioStore) scanAll(rows pgx.Rows) (*[]gaivota.Portfolio, error) {
	var portfolios []gaivota.Portfolio

	for rows.Next() {
		portfolio, err := store.scanOne(rows)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning portfolios: %w", err)
		}

		portfolios = append(portfolios, *portfolio)
	}

	return &portfolios, nil
}

func (store *PortfolioStore) scanOne(row pgx.Row) (*gaivota.Portfolio, error) {
	var portfolio gaivota.Portfolio

	err := row.Scan(&portfolio.ID, &portfolio.UserID, &portfolio.Name, &portfolio.CreatedAt, &portfolio.UpdatedAt, &portfolio.DeletedAt)

	return &portfolio, err
}

func (store *PortfolioStore) Add(ctx context.Context, portfolio *gaivota.Portfolio) (*gaivota.Portfolio, error) {
	query := `insert into portfolios ("user_id", "name")
						values ($1, $2)
						returning "id", "user_id", "name", "created_at", "updated_at", "deleted_at"`

	row := store.Database.Pool.QueryRow(ctx, query, portfolio.UserID, portfolio.Name)

	newPortfolio, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not insert portfolio %s for user %v: %w", portfolio.Name, portfolio.UserID, err)
	}

	return newPortfolio, nil
}

func (store *PortfolioStore) All(ctx context.Context) (*[]gaivota.Portfolio, error) {
	query := `select "id", "user_id", "name", "created_at", "updated_at", "deleted_at"
						from portfolios`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get portfolios: %w", err)
	}

	return store.scanAll(rows)
}

func (store *PortfolioStore) Delete(ctx context.Context, id int) error {
	query := `update portfolios
						set deleted_at = now(),
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, id,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete portfolio %v: %w", id, err)
	}

	return nil
}

func (store *PortfolioStore) Get(ctx context.Context, id int) (*gaivota.Portfolio, error) {
	query := `select "id", "user_id", "name", "created_at", "updated_at", "deleted_at"
						from portfolios where id = $1`

	portfolio := &gaivota.Portfolio{}

	row := store.Database.Pool.QueryRow(ctx, query, id)

	portfolio, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not get portfolio %v: %w", id, err)
	}

	return portfolio, nil
}

func (store *PortfolioStore) GetByUserID(ctx context.Context, userId int) (*[]gaivota.Portfolio, error) {
	query := `select "id", "user_id", "name", "created_at", "updated_at", "deleted_at"
						from portfolios where user_id = $1`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get portfolios for user %v: %w", userId, err)
	}

	return store.scanAll(rows)
}

func (store *PortfolioStore) Update(ctx context.Context, portfolio *gaivota.Portfolio) error {
	query := `update portfolios
						set name = $1
						where id = $2`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, &portfolio.Name, &portfolio.ID,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update portfolio %v: %w", portfolio.ID, err)
	}

	return nil
}
