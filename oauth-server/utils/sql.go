package utils

import (
	"context"
	"fmt"
	custom_types "local/bomboclat-oauth-server/types"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	pgInstance *custom_types.Postgres
	pgOnce     sync.Once
)

func CreateDBConnPool() (*custom_types.Postgres, error) {

	godotenv.Load()

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

func CreateClientsTable(pg *custom_types.Postgres) error {
	query := `CREATE TABLE IF NOT EXISTS clients (
		client_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		client_secret_hash TEXT,
		redirect_uri TEXT,
		app_name TEXT,
		logo_url TEXT,
		grant_types JSONB,
		jwks_uri TEXT,
		is_confidential BOOLEAN,
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	)`

	_, err := pg.DB.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}

// Part of OIDC, not OAuth
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

//code	VARCHAR(255)	Random unique auth code
//user_id	FK → users.id	Who authorized
//client_id	FK → clients.id	Which client
//redirect_uri	TEXT	For verification in /token
//scopes	TEXT	Space-separated scope
//expires_at	TIMESTAMP	Expiry (30–120 sec)
//code_challenge	TEXT	PKCE support (nullable)
//code_challenge_method	VARCHAR(10)	e.g., "S256"
//used	BOOLEAN	One-time use flag
//created_at	TIMESTAMP	Timestamp

func CreateAuthCodesTable(pg *custom_types.Postgres) error {
	query := `CREATE TABLE IF NOT EXISTS auth_codes (
		id SERIAL PRIMARY KEY,
		code TEXT,
		user_id UUID,
		client_id UUID,
		redirect_uri TEXT,
		scopes TEXT,
		expires_at TIMESTAMP DEFAULT now() + INTERVAL '10 minutes',
		code_challenge TEXT,
		code_challenge_method VARCHAR(10),
		used BOOLEAN,
		created_at TIMESTAMP DEFAULT now(),
		FOREIGN KEY (client_id) REFERENCES clients(client_id),
		FOREIGN KEY (user_id) REFERENCES users(uuid)
	)`

	_, err := pg.DB.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}

//token	TEXT	The actual access token (JWT or random string)
//user_id	FK → users.id	Token owner
//client_id	FK → clients.id	Which client owns this token
//scopes	TEXT	Permissions granted
//expires_at	TIMESTAMP	Usually 1 hour
//created_at	TIMESTAMP	Issued time
//revoked	BOOLEAN	Revocation support

func CreateAccessTokensTable(pg *custom_types.Postgres) error {
	query := `CREATE TABLE IF NOT EXISTS access_tokens (
		token TEXT,
		user_id INTEGER,
		client_id UUID,
		scopes TEXT,
		expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT now(),
		revoked BOOLEAN,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (client_id) REFERENCES clients(client_id)
	)`

	_, err := pg.DB.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}

func CreateRefreshTokensTable(pg *custom_types.Postgres) error {
	query := `CREATE TABLE IF NOT EXISTS refresh_tokens (
		token TEXT,
		user_id INTEGER,
		client_id UUID,
		scopes TEXT,
		expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT now(),
		revoked BOOLEAN,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (client_id) REFERENCES clients(client_id)
	)`

	_, err := pg.DB.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}

func CreateConsentsTable(pg *custom_types.Postgres) error {
	query := `CREATE TABLE IF NOT EXISTS consents (
		id SERIAL PRIMARY KEY,
		user_id INTEGER,
		client_id UUID,
		scopes TEXT,
		updated_at TIMESTAMP DEFAULT now(),
		created_at TIMESTAMP DEFAULT now(),
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (client_id) REFERENCES clients(client_id)
	)`

	_, err := pg.DB.Exec(context.Background(), query)

	if err != nil {
		return err
	}

	return nil
}
