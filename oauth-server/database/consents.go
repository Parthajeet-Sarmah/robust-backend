package database

import (
	"context"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/jackc/pgx/v5"
)

// FindConsent checks if a user has previously granted consent to a client
func FindConsent(pg *custom_types.Postgres, userId string, clientId string) (*models.ConsentModel, error) {
	query := `SELECT * FROM consents WHERE user_id = @userId AND client_id = @clientId LIMIT 1`
	args := pgx.NamedArgs{
		"userId":   userId,
		"clientId": clientId,
	}

	rows, err := pg.DB.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.ConsentModel])
	if err != nil {
		// If no rows found, return nil without error (user hasn't consented yet)
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

// UpsertConsent creates or updates a consent record
func UpsertConsent(pg *custom_types.Postgres, m *models.ConsentModel) error {
	query := `INSERT INTO consents (user_id, client_id, scopes, updated_at, created_at)
		VALUES (@userId, @clientId, @scopes, NOW(), NOW())
		ON CONFLICT (user_id, client_id)
		DO UPDATE SET scopes = @scopes, updated_at = NOW()`

	args := pgx.NamedArgs{
		"userId":   m.UserId,
		"clientId": m.ClientId,
		"scopes":   m.Scopes,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)
	return err
}

// DeleteConsent removes a consent record (for revocation)
func DeleteConsent(pg *custom_types.Postgres, userId string, clientId string) error {
	query := `DELETE FROM consents WHERE user_id = @userId AND client_id = @clientId`
	args := pgx.NamedArgs{
		"userId":   userId,
		"clientId": clientId,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)
	return err
}
