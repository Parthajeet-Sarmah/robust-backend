package database

import (
	"context"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/jackc/pgx/v5"
)

func FindUserByEmailAndPasswordHash(pg *custom_types.Postgres, ctx context.Context, email string, password_hash string) (*models.UserDatabaseModel, error) {

	query := `SELECT * FROM users WHERE email = @email AND password_hash = @pHash LIMIT 1`
	args := pgx.NamedArgs{"email": email, "pHash": password_hash}
	rows, err := pg.DB.Query(ctx, query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserDatabaseModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func InsertUser(pg *custom_types.Postgres, ctx context.Context, m *custom_types.UserDatabaseModelInput) error {
	query := `INSERT INTO users (
		username,
		email,
		password_hash
	) VALUES (
		@userName,
		@email,
		@passwordHash
	)`

	args := pgx.NamedArgs{
		"userName":     m.Username,
		"email":        m.Email,
		"passwordHash": m.PasswordHash,
	}

	_, err := pg.DB.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}

//	query := `CREATE TABLE IF NOT EXISTS consents (
//		id SERIAL PRIMARY KEY,
//		user_id UUID,
//		client_id UUID,
//		scopes TEXT,
//		updated_at TIMESTAMP DEFAULT now(),
//		created_at TIMESTAMP DEFAULT now(),
//		FOREIGN KEY (client_id) REFERENCES clients(client_id)
//	)`

func FindUserConsent(pg *custom_types.Postgres, ctx context.Context, userId string, clientId string) (*models.UserConsentsModel, error) {
	query := `SELECT user_id, client_id, scopes FROM consents WHERE user_id = @userId AND client_id = @clientId LIMIT 1`
	args := pgx.NamedArgs{"userId": userId, "clientId": clientId}
	rows, err := pg.DB.Query(ctx, query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.UserConsentsModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func InsertUserConsent(pg *custom_types.Postgres, ctx context.Context, m *custom_types.UserConsentInput) error {
	query := `INSERT INTO consents (
		user_id,
		client_id,
		scopes
	) VALUES (
		@userId,
		@clientId,
		@scopes
	)`

	args := pgx.NamedArgs{
		"userId":   m.UserId,
		"clientId": m.ClientId,
		"scopes":   m.Scopes,
	}

	_, err := pg.DB.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}
