package database

import (
	"context"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/jackc/pgx/v5"
)

func InsertAccessToken(pg *custom_types.Postgres, m *models.AccessTokenModel) error {
	query := `INSERT INTO access_tokens (
		token_hash,
		user_id,
		client_id,
		scopes,
		expires_at
	) VALUES (
		@tokenHash,
		@userId,
		@clientId,
		@scopes,
		@expiresAt
	)`

	args := pgx.NamedArgs{
		"tokenHash": m.TokenHash,
		"userId":    m.UserId,
		"clientId":  m.ClientId,
		"scopes":    m.Scopes,
		"expiresAt": m.ExpiresAt,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil
}

func RevokeAccessToken(pg *custom_types.Postgres, tokenHash string) error {

	query := `UPDATE access_tokens SET revoked = true WHERE token_hash = @tokenHash`
	args := pgx.NamedArgs{
		"tokenHash": tokenHash,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil

}

func UpdateAccessToken(pg *custom_types.Postgres, m *models.AccessTokenModel) error {

	query := `UPDATE access_tokens SET
		token_hash = @tokenHash,
		expires_at = @expiresAt,
		scopes = @scopes
	WHERE user_id = @userId AND client_id = @clientId`

	args := pgx.NamedArgs{
		"tokenHash": m.TokenHash,
		"userId":    m.UserId,
		"clientId":  m.ClientId,
		"scopes":    m.Scopes,
		"expiresAt": m.ExpiresAt,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil

}

func InsertRefreshToken(pg *custom_types.Postgres, m *models.RefreshTokenModel) error {

	query := `INSERT INTO refresh_tokens (
		token_hash,
		user_id,
		client_id,
		scopes,
		expires_at
	) VALUES (
		@tokenHash,
		@userId,
		@clientId,
		@scopes,
		@expiresAt
	)`

	args := pgx.NamedArgs{
		"tokenHash": m.TokenHash,
		"userId":    m.UserId,
		"clientId":  m.ClientId,
		"scopes":    m.Scopes,
		"expiresAt": m.ExpiresAt,
	}

	_, err := pg.DB.Exec(context.Background(), query, args)

	if err != nil {
		return err
	}

	return nil

}

func FindAccessToken(pg *custom_types.Postgres, accessTokenHash string) (*models.AccessTokenModel, error) {

	query := `SELECT * FROM access_tokens WHERE token_hash = @tokenHash LIMIT 1`
	args := pgx.NamedArgs{
		"tokenHash": accessTokenHash,
	}

	rows, err := pg.DB.Query(context.Background(), query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.AccessTokenModel])

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func IntrospectAccessToken(pg *custom_types.Postgres, accessTokenHash string) (*models.IntrospectAccessTokenModel, error) {

	query := `SELECT client_id, scopes, revoked FROM access_tokens WHERE token_hash = @tokenHash LIMIT 1`
	args := pgx.NamedArgs{
		"tokenHash": accessTokenHash,
	}

	rows, err := pg.DB.Query(context.Background(), query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.IntrospectAccessTokenModel])

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func FindRefreshToken(pg *custom_types.Postgres, refreshTokenHash string) (*models.RefreshTokenModel, error) {

	query := `SELECT * FROM refresh_tokens WHERE token_hash = @tokenHash LIMIT 1`
	args := pgx.NamedArgs{
		"tokenHash": refreshTokenHash,
	}

	rows, err := pg.DB.Query(context.Background(), query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.RefreshTokenModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func IntrospectRefreshToken(pg *custom_types.Postgres, refreshTokenHash string) (*models.IntrospectRefreshTokenModel, error) {

	query := `SELECT client_id, expires_at, created_at, revoked, scopes FROM refresh_tokens WHERE token_hash = @tokenHash LIMIT 1`
	args := pgx.NamedArgs{
		"tokenHash": refreshTokenHash,
	}

	rows, err := pg.DB.Query(context.Background(), query, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.IntrospectRefreshTokenModel])

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func UpdateRefreshTokenEntry(pg *custom_types.Postgres, oldTokenHash string, newTokenHash string) error {

	query := `UPDATE refresh_tokens SET token_hash = @newTokenHash WHERE token_hash = @oldTokenHash`
	args := pgx.NamedArgs{
		"oldTokenHash": oldTokenHash,
		"newTokenHash": newTokenHash,
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
