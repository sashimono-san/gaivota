package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/leoschet/gaivota"
)

type HoldingStore struct {
	Database *Database
}

func (store *HoldingStore) getByFK(ctx context.Context, fk_column string, fk int) (*[]gaivota.Holding, error) {
	query := `select "id", "wallet_id", "position_id", "amount", "created_at", "updated_at", "deleted_at"
						from holdings where $1 = $2`

	rows, err := store.Database.Pool.Query(ctx, query, fk_column, fk)

	if err != nil {
		return nil, fmt.Errorf("Could not get holdings where %s is %v: %w", fk_column, fk, err)
	}

	return store.scanAll(rows)
}

func (store *HoldingStore) scanAll(rows pgx.Rows) (*[]gaivota.Holding, error) {
	var holdings []gaivota.Holding

	for rows.Next() {
		holding, err := store.scanOne(rows)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning holdings: %w", err)
		}

		holdings = append(holdings, *holding)
	}

	return &holdings, nil
}

func (store *HoldingStore) scanOne(row pgx.Row) (*gaivota.Holding, error) {
	var holding gaivota.Holding

	err := row.Scan(
		&holding.ID, &holding.WalletID, &holding.PositionID, &holding.Amount,
		&holding.CreatedAt, &holding.UpdatedAt, &holding.DeletedAt,
	)

	return &holding, err
}

func (store *HoldingStore) Add(ctx context.Context, holding *gaivota.Holding) (*gaivota.Holding, error) {
	query := `insert into holdings ("wallet_id", "position_id", "amount")
						values ($1, $2)
						returning "id", "wallet_id", "position_id", "amount", "created_at", "updated_at", "deleted_at"`

	row := store.Database.Pool.QueryRow(
		ctx, query, holding.WalletID, holding.PositionID, holding.Amount,
	)

	newHolding, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf(
			"Could not insert holding for wallet %v and position %v: %w",
			holding.WalletID, holding.PositionID, err,
		)
	}

	return newHolding, nil
}

func (store *HoldingStore) All(ctx context.Context) (*[]gaivota.Holding, error) {
	query := `select "id", "wallet_id", "position_id", "amount", "created_at", "updated_at", "deleted_at"
						from holdings`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get holdings: %w", err)
	}

	return store.scanAll(rows)
}

func (store *HoldingStore) Delete(ctx context.Context, id int) error {
	query := `update holdings
						set deleted_at = now(),
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(ctx, query, id)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete holding %v: %w", id, err)
	}

	return nil
}

func (store *HoldingStore) Get(ctx context.Context, id int) (*gaivota.Holding, error) {
	query := `select "id", "wallet_id", "position_id", "amount", "created_at", "updated_at", "deleted_at"
						from holdings where id = $1`

	row := store.Database.Pool.QueryRow(
		ctx, query, id,
	)

	holding, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not get holding %v: %w", id, err)
	}

	return holding, nil
}

func (store *HoldingStore) GetByUserID(ctx context.Context, userId int) (*[]gaivota.Holding, error) {
	query := `select "h.id", "h.wallet_id", "h.position_id", "h.amount", "h.created_at", "h.updated_at", "h.deleted_at"
						from holdings as h
						join wallets as w on "w.id" = "h.wallet_id"
						where w.user_id = $1`

	rows := store.Database.Pool.Query(ctx, query, userId)

	return store.scanAll(rows)
}

func (store *HoldingStore) GetByWalletID(ctx context.Context, walletId int) (*[]gaivota.Holding, error) {
	return store.getByFK(ctx, "wallet_id", walletId)
}

func (store *HoldingStore) GetByPositionID(ctx context.Context, positionId int) (*[]gaivota.Holding, error) {
	return store.getByFK(ctx, "position_id", positionId)
}

func (store *HoldingStore) Update(ctx context.Context, holding *gaivota.Holding) error {
	query := `update holdings
						set wallet_id = $1,
								position_id = $2,
								amount = $3,
						where id = $4`

	cmdTags, err := store.Database.Pool.Exec(ctx, query, &holding.WalletID, &holding.PositionID, &holding.Amount, &holding.ID)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update holding %v: %w", holding.ID, err)
	}

	return nil
}
