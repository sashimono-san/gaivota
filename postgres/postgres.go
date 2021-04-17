package postgres

import (
	"context"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

type Database struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, connString string) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// TODO: Add a logger
	// poolConfig.ConnConfig.Logger = logger

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	db := &Database{Pool: pool}
	return db, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}

func (db *Database) Ping() (msg string, err error) {
	err = db.Pool.Ping(context.Background())

	if err != nil {
		return "Could not connect to the Database", err
	}

	return "", nil
}

// func (db *Database) CreateTable(model interface{}) {
// 	// Get how many properties are in the model
// 	column_count := reflect.ValueOf(model).NumField()

// 	columns_name := make([]interface{}, column_count)

// 	for i := 0; i < column_count; i++ {
// 		// Get the name of each property in snake_case format
// 		columns_name[i] = toSnakeCase(reflect.TypeOf(model).Field(i).Name)
// 	}

// 	fmt.Println(columns_name)
// }
