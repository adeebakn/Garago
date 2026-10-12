package app

import (
	"context"

	db "garago/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// databasePool provides independent connections for concurrent API requests.
// A pool is required because one pgx connection cannot process two queries at
// the same time (which caused the previous "conn busy" error).
var databasePool *pgxpool.Pool

func connectDB(ctx context.Context, databaseURL string) (*db.Queries, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)

	if err != nil {
		return nil, nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, err
	}

	queries := db.New(pool)
	databasePool = pool

	return queries, pool, nil
}
