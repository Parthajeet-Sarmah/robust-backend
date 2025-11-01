package database

import (
	"context"
	"fmt"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/jackc/pgx/v5"
)

func FindClientById(pg *custom_types.Postgres, ctx context.Context, client_id string) (*models.ClientDatabaseModel, error) {

	query := `SELECT * FROM clients WHERE client_id = @clientId LIMIT 1`
	args := pgx.NamedArgs{"clientId": client_id}
	rows, err := pg.DB.Query(ctx, query, args)

	if err != nil {
		return nil, err
	}

	fmt.Println(client_id)

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ClientDatabaseModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func InsertClient(pg *custom_types.Postgres, ctx context.Context, m *custom_types.ClientDatabaseModelInput) error {
	query := `INSERT INTO clients (
		client_secret_hash,
		redirect_uri,
		app_name,
		logo_url,
		grant_types,
		jwks_uri,
		is_confidential
	) VALUES (
		@clientSecretHash,
		@redirectUri,
		@appName,
		@logoUrl,
		@grantTypes,
		@jwksUri,
		@isConfidential
	)`

	args := pgx.NamedArgs{
		"clientSecretHash": m.ClientSecretHash,
		"redirectUri":      m.RedirectUri,
		"appName":          m.AppName,
		"logoUrl":          m.LogoUrl,
		"grantTypes":       m.GrantTypes,
		"jwksUri":          m.JwksUri,
		"isConfidential":   m.IsConfidential,
	}

	_, err := pg.DB.Exec(ctx, query, args)

	if err != nil {
		return err
	}

	return nil
}
