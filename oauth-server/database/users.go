package database

import (
	"context"
	"fmt"
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

	fmt.Println(email)

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.UserDatabaseModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func InsertUser(pg *custom_types.Postgres, ctx context.Context, m *models.UserDatabaseModelInput) error {
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
