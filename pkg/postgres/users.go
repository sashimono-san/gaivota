package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgx"
	"github.com/leoschet/gaivota"
)

type UserStore struct {
	Database *Database
}

func (store *UserStore) scanAll(rows pgx.Rows) (*[]gaivota.User, error) {
	var users []gaivota.User

	for rows.Next() {
		user, err := store.scanOne(rows)

		if err != nil {
			return nil, fmt.Errorf("Error while scanning users: %w", err)
		}

		users = append(users, *user)
	}

	return &users, nil
}

func (store *UserStore) scanOne(row pgx.Row) (*gaivota.User, error) {
	var user gaivota.User

	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	return &user, err
}

func (store *UserStore) Add(ctx context.Context, user *gaivota.User) (*gaivota.User, error) {
	query := `insert into users ("email", "first_name", "last_name")
						values ($1, $2, $3)
						returning "id", "email", "first_name", "last_name", "created_at", "updated_at", "deleted_at"`

	row := store.Database.Pool.QueryRow(ctx, query, user.Email, user.FirstName, user.LastName)

	newUser, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not insert user %s: %w", user.Email, err)
	}

	return newUser, nil
}

func (store *UserStore) All(ctx context.Context) (*[]gaivota.User, error) {
	query := `select "id", "email", "first_name", "last_name", "created_at", "updated_at", "deleted_at"
						from users`

	rows, err := store.Database.Pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Could not get users: %w", err)
	}

	return store.scanAll(rows)
}

func (store *UserStore) Delete(ctx context.Context, id int) error {
	query := `update users
						set deleted_at = now(),
						where id = $1`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, id,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not delete user %v: %w", id, err)
	}

	return nil
}

func (store *UserStore) Get(ctx context.Context, id int) (*gaivota.User, error) {
	query := `select "id", "email", "first_name", "last_name", "created_at", "updated_at", "deleted_at"
						from users where id = $1`

	row := store.Database.Pool.QueryRow(ctx, query, id)

	user, err := store.scanOne(row)

	if err != nil {
		return nil, fmt.Errorf("Could not get user %v: %w", id, err)
	}

	return user, nil
}

func (store *UserStore) Update(ctx context.Context, user *gaivota.User) error {
	query := `update users
						set email = $1,
								first_name = $2,
								last_name = $3
						where id = $4`

	cmdTags, err := store.Database.Pool.Exec(
		ctx, query, user.Email, user.FirstName, user.LastName, user.ID,
	)

	if err != nil || cmdTags.RowsAffected() == 0 {
		return fmt.Errorf("Could not update user %v: %w", user.ID, err)
	}

	return nil
}
