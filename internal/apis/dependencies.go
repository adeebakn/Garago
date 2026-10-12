package apis

import "github.com/jackc/pgx/v5/pgxpool"

// databasePool is shared by handlers that require explicit transactions or
// SQL outside the generated query surface. New feature code should prefer
// generated queries and receive dependencies through constructors.
var databasePool *pgxpool.Pool
