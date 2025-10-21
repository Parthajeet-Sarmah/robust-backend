package custom_types

import "github.com/jackc/pgx/v5/pgxpool"

type UserDetails struct {
	Email        string
	PasswordHash string
}

type Postgres struct {
	DB *pgxpool.Pool
}
