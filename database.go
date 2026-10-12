package main

import (
	"context"
	"os"

	db "garago/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// databasePool provides independent connections for concurrent API requests.
// A pool is required because one pgx connection cannot process two queries at
// the same time (which caused the previous "conn busy" error).
var databasePool *pgxpool.Pool

func connectDB() (*db.Queries, *pgxpool.Pool, error) {

	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		return nil, nil, err
	}

	queries := db.New(pool)
	databasePool = pool

	return queries, pool, nil
}
