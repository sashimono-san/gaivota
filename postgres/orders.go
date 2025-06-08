package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/leoschet/gaivota"
)

func NewOrderStore(db *Database) *OrderStore {
	return &OrderStore{
		Database: db,
	}
}

type OrderStore struct {
	Database *Database
}

func (store *OrderStore) scanAll(rows pgx.Rows) ([]gaivota.Order, error) {
	var orders []gaivota.Order

	for rows.Next() {
		order, err := store.scanOne(rows)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning orders: %w", err)
		}

		orders = append(orders, *order)
	}

	return orders, nil
}

func (store *OrderStore) scanOne(row pgx.Row) (*gaivota.Order, error) {
	var order gaivota.Order

	err := row.Scan(
		&order.ID, &order.PositionID, &order.Amount, &order.UnitPrice, &order.TotalPrice,
		&order.Operation, &order.Type, &order.Exchange, &order.ExecutedAt,
		&order.CreatedAt, &order.UpdatedAt, &order.DeletedAt,
	)

	return &order, err
}

func (store *OrderStore) Add(ctx context.Context, order *gaivota.Order) (*gaivota.Order, error) {
	query := `insert into orders ("position_id", "amount", "unit_price", "total_price", "operation", "type", "exchange", "executed_at")
						values ($1, $2, $3, $4, $5, $6, $7, $8)
						returning "id", "position_id", "amount", "unit_price", "total_price", "operation", "type", "exchange", "executed_at", "created_at", "updated_at", "deleted_at"`

	row := store.Database.Pool.QueryRow(
		ctx, query, order.PositionID, order.Amount, order.UnitPrice, order.TotalPrice,
		order.Operation, order.Type, order.Exchange, order.ExecutedAt,
	)

	newOrder, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf(
			"Could not insert order for position %v: %w",
			order.PositionID, err,
		)
	}

	return newOrder, nil
}

func (store *OrderStore) All(ctx context.Context) ([]gaivota.Order, error) {
	query := `select "id", "position_id", "amount", "unit_price", "total_price", "operation", "type", "exchange", "executed_at", "created_at", "updated_at", "deleted_at"
						from orders where deleted_at is null`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get orders: %w", err)
	}

	return store.scanAll(rows)
}

func (store *OrderStore) Delete(ctx context.Context, id int) error {
	query := `update orders
						set deleted_at = now()
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(ctx, query, id)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete order %v: %w", id, err)
	}

	return nil
}

func (store *OrderStore) Get(ctx context.Context, id int) (*gaivota.Order, error) {
	query := `select "id", "position_id", "amount", "unit_price", "total_price", "operation", "type", "exchange", "executed_at", "created_at", "updated_at", "deleted_at"
						from orders where id = $1 and deleted_at is null`

	row := store.Database.Pool.QueryRow(ctx, query, id)

	order, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not get order %v: %w", id, err)
	}

	return order, nil
}

func (store *OrderStore) Update(ctx context.Context, order *gaivota.Order) error {
	query := `update orders
						set position_id = $1,
								amount = $2,
								unit_price = $3,
								total_price = $4,
								operation = $5,
								type = $6,
								exchange = $7,
								executed_at = $8
						where id = $9`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, order.PositionID, order.Amount, order.UnitPrice, order.TotalPrice,
		order.Operation, order.Type, order.Exchange, order.ExecutedAt, order.ID,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update order %v: %w", order.ID, err)
	}

	return nil
}
