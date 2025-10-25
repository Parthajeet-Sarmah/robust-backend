package database

import (
	"context"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/jackc/pgx/v5"
)

func InsertAccessToken(pg *custom_types.Postgres, m *models.AccessTokenModel) error {
	query := `INSERT INTO access_tokens (
		token,
		user_id,
		client_id,
		scopes
	) VALUES (
		@token,
		@userId,
		@clientId,
		@scopes
	)`

	args := pgx.NamedArgs{
		"token":    m.Token,
		"userId":   m.UserId,
		"clientId": m.ClientId,
		"scopes":   m.Scopes,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

func InsertRefreshToken(pg *custom_types.Postgres, m *models.RefreshTokenModel) error {

	query := `INSERT INTO refresh_tokens (
		token,
		user_id,
		client_id,
		scopes
	) VALUES (
		@token,
		@userId,
		@clientId,
		@scopes
	)`

	args := pgx.NamedArgs{
		"token":    m.Token,
		"userId":   m.UserId,
		"clientId": m.ClientId,
		"scopes":   m.Scopes,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil

}

func UpdateAuthCodeEntryUsedStatus(pg *custom_types.Postgres, code string) error {

	query := `UPDATE auth_codes SET used = @used WHERE code = @code`
	args := pgx.NamedArgs{
		"code": code,
		"used": true,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

func GetAuthCode(pg *custom_types.Postgres, code string) (*models.AuthCodeModel, error) {

	query := `SELECT * FROM auth_codes WHERE code = @code LIMIT 1`
	args := pgx.NamedArgs{
		"code": code,
	}
	rows, err := pg.DB.Query(context.Background(), query, args)

	if err != nil {
		return nil, err

	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.AuthCodeModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func CreateAuthCodeEntry(pg *custom_types.Postgres, m *models.AuthCodeModelInput) error {

	query := `INSERT INTO auth_codes (code, user_id, client_id, redirect_uri, scopes, code_challenge,
	code_challenge_method, used) VALUES (
		@code,
		@userId,
		@clientId,
		@redirectUri,
		@scopes,
		@codeChallenge,
		@codeChallengeMethod,
		@used
	)`

	args := pgx.NamedArgs{
		"code":                m.Code,
		"userId":              m.UserId,
		"clientId":            m.ClientId,
		"redirectUri":         m.RedirectUri,
		"scopes":              m.Scopes,
		"codeChallenge":       m.CodeChallenge,
		"codeChallengeMethod": m.CodeChallengeMethod,
		"used":                m.Used,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}
