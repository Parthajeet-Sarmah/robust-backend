package utils

import (
	"context"
	"fmt"
	custom_types "local/bomboclat-oidc-service/types"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *custom_types.Postgres
	pgOnce     sync.Once
)

func CreateDBConnPool() (*custom_types.Postgres, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")
	user := os.Getenv("DB_USER")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, pass, host, port, name)

	pgOnce.Do(func() {
		db, err := pgxpool.New(context.Background(), connStr)

		if err != nil {
			log.Fatal(err)
			return
		}

		pgInstance = &custom_types.Postgres{
			DB: db,
		}
	})

	return pgInstance, nil
}
func CreateUsersTable(pg *custom_types.Postgres) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT,
		email TEXT,
		password_hash TEXT,
		uuid UUID UNIQUE DEFAULT gen_random_uuid(),
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	)`

	_, err := pg.DB.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}
