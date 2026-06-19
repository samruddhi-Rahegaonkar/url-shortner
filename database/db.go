package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	var err error

	DB, err = pgxpool.New(context.Background(),
		"postgres://postgres:root@localhost:5432/urlshortner",
	)
	if err != nil {
		panic(err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connected successfully")

}
