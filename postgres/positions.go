package postgres

import (
	"context"
	"fmt"

	"github.com/leoschet/gaivota"
)

type PositionStore struct {
	Database *Database
}

func (store *PositionStore) Add(ctx context.Context, position *gaivota.Position) (*gaivota.Position, error) {
	query := `insert into positions ("investment_id", "amount", "average_price", "profit")
						values ($1, $2)
						returning "id", "investment_id", "amount", "average_price", "profit", "created_at", "updated_at", "deleted_at"`

	var newPosition gaivota.Position

	err := store.Database.Pool.QueryRow(
		ctx, query, position.InvestmentID, position.Amount,
		position.AveragePrice, position.Profit,
	).Scan(
		&newPosition.ID, &newPosition.InvestmentID, &newPosition.Amount,
		&newPosition.AveragePrice, &newPosition.Profit,
		&newPosition.CreatedAt, &newPosition.UpdatedAt, &newPosition.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not insert position for investment %v: %w", position.InvestmentID, err)
	}

	return &newPosition, nil
}

func (store *PositionStore) All(ctx context.Context) (*[]gaivota.Position, error) {
	query := `select "id", "investment_id", "amount", "average_price", "profit", "created_at", "updated_at", "deleted_at"
						from positions`

	var positions []gaivota.Position
	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get positions: %w", err)
	}

	for rows.Next() {
		var position gaivota.Position
		err = rows.Scan(
			&position.ID, &position.InvestmentID, &position.Amount,
			&position.AveragePrice, &position.Profit,
			&position.CreatedAt, &position.UpdatedAt, &position.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning positions: %w", err)
		}

		positions = append(positions, position)
	}

	return &positions, nil
}

func (store *PositionStore) Delete(ctx context.Context, id int) error {
	query := `update positions
						set deleted_at = now(),
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, id,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete position %v: %w", id, err)
	}

	return nil
}

func (store *PositionStore) Get(ctx context.Context, id int) (*gaivota.Position, error) {
	query := `select "id", "investment_id", "amount", "average_price", "profit", "created_at", "updated_at", "deleted_at"
						from positions where id = $1`

	position := &gaivota.Position{}

	err := store.Database.Pool.QueryRow(
		ctx, query, id,
	).Scan(
		&position.ID, &position.InvestmentID, &position.Amount,
		&position.AveragePrice, &position.Profit,
		&position.CreatedAt, &position.UpdatedAt, &position.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not get position %v: %w", id, err)
	}

	return position, nil
}

func (store *PositionStore) GetByUserID(ctx context.Context, userId int) (*gaivota.Position, error) {
	query := `select "id", "investment_id", "amount", "average_price", "profit", "created_at", "updated_at", "deleted_at"
						from positions where investment_id = $1`

	position := &gaivota.Position{}

	err := store.Database.Pool.QueryRow(
		ctx, query, userId,
	).Scan(
		&position.ID, &position.UserID, &position.Name,
		&position.TotalValue, &position.Address, &position.Location,
		&position.CreatedAt, &position.UpdatedAt, &position.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("Could not get positions for user %v: %w", userId, err)
	}

	return position, nil
}

func (store *PositionStore) Update(ctx context.Context, position *gaivota.Position) error {
	query := `update positions
						set name = $1,
								total_value = $2
						where id = $5`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, &position.Name, &position.TotalValue, &position.ID,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update position %v: %w", position.ID, err)
	}

	return nil
}
