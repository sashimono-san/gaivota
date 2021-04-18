package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/leoschet/gaivota"
)

func NewWalletStore(db *Database) WalletStore {
	return WalletStore{
		Database: db,
	}
}

type WalletStore struct {
	Database *Database
}

func (store *WalletStore) scanAll(rows pgx.Rows) (*[]gaivota.Wallet, error) {
	var wallets []gaivota.Wallet

	for rows.Next() {
		wallet, err := store.scanOne(rows)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning wallets: %w", err)
		}

		wallets = append(wallets, *wallet)
	}

	return &wallets, nil
}

func (store *WalletStore) scanOne(row pgx.Row) (*gaivota.Wallet, error) {
	var wallet gaivota.Wallet

	err := row.Scan(
		&wallet.ID, &wallet.UserID, &wallet.Name,
		&wallet.TotalValue, &wallet.Address, &wallet.Location,
		&wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt,
	)

	return &wallet, err
}

func (store *WalletStore) Add(ctx context.Context, wallet *gaivota.Wallet) (*gaivota.Wallet, error) {
	query := `insert into wallets ("user_id", "name", "total_value", "address", "location")
						values ($1, $2)
						returning "id", "user_id", "name", "total_value", "address", "location", "created_at", "updated_at", "deleted_at"`

	row := store.Database.Pool.QueryRow(
		ctx, query, wallet.UserID, wallet.Name,
	)

	newWallet, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not insert wallet %s for user %v: %w", wallet.Name, wallet.UserID, err)
	}

	return newWallet, nil
}

func (store *WalletStore) All(ctx context.Context) (*[]gaivota.Wallet, error) {
	query := `select "id", "user_id", "name", "total_value", "address", "location", "created_at", "updated_at", "deleted_at"
						from wallets`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get wallets: %w", err)
	}

	return store.scanAll(rows)
}

func (store *WalletStore) Delete(ctx context.Context, id int) error {
	query := `update wallets
						set deleted_at = now(),
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, id,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete wallet %v: %w", id, err)
	}

	return nil
}

func (store *WalletStore) Get(ctx context.Context, id int) (*gaivota.Wallet, error) {
	query := `select "id", "user_id", "name", "total_value", "address", "location", "created_at", "updated_at", "deleted_at"
						from wallets where id = $1`

	row := store.Database.Pool.QueryRow(
		ctx, query, id,
	)

	wallet, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not get wallet %v: %w", id, err)
	}

	return wallet, nil
}

func (store *WalletStore) GetByUserID(ctx context.Context, userId int) (*[]gaivota.Wallet, error) {
	query := `select "id", "user_id", "name", "total_value", "address", "location", "created_at", "updated_at", "deleted_at"
						from wallets where user_id = $1`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get wallets for user %v: %w", userId, err)
	}

	return store.scanAll(rows)
}

func (store *WalletStore) Update(ctx context.Context, wallet *gaivota.Wallet) error {
	query := `update wallets
						set name = $1,
								total_value = $2
						where id = $3`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, &wallet.Name, &wallet.TotalValue, &wallet.ID,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update wallet %v: %w", wallet.ID, err)
	}

	return nil
}
