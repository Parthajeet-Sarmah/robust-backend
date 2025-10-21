package database

import (
	"context"
	"local/bomboclat-oauth-server/models"
	custom_types "local/bomboclat-oauth-server/types"

	"github.com/jackc/pgx/v5"
)

func CreateAuthCodeEntry(pg *custom_types.Postgres, m models.AuthCodeModelInput) error {

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
