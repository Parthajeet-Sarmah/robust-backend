package database

import (
	"context"
	"local/bomboclat-oidc-service/models"
	custom_types "local/bomboclat-oidc-service/types"
	"local/bomboclat-oidc-service/utils"

	"github.com/jackc/pgx/v5"
)

func FindUserByUUID(pg *custom_types.Postgres, ctx context.Context, uuid string) (*models.UserDatabaseModel, error) {

	query := `SELECT username, email FROM users WHERE uuid = @uuid LIMIT 1`
	args := pgx.NamedArgs{"uuid": uuid}
	rows, err := pg.DB.Query(ctx, query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.UserDatabaseModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func FindUserByEmail(pg *custom_types.Postgres, ctx context.Context, email string) (*models.UserDatabaseModel, error) {

	query := `SELECT * FROM users WHERE email = @email LIMIT 1`
	args := pgx.NamedArgs{"email": email}
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

func InsertUser(pg *custom_types.Postgres, m *custom_types.UserRegistrationDetails) error {
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
		"passwordHash": utils.HashToken256(m.Password),
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}
