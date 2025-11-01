package models

import (
	"time"
)

type AuthCodeModel struct {
	Id                  string    `db:"id"`
	Code                string    `db:"code"`
	UserId              string    `db:"user_id"`
	ClientId            string    `db:"client_id"`
	RedirectUri         string    `db:"redirect_uri"`
	Scopes              string    `db:"scopes"`
	ExpiresAt           time.Time `db:"expires_at"`
	CodeChallenge       string    `db:"code_challenge"`
	CodeChallengeMethod string    `db:"code_challenge_method"`
	Used                bool      `db:"used"`
	CreatedAt           time.Time `db:"created_at"`
}

type AccessTokenModel struct {
	TokenHash string
	UserId    string
	ClientId  string
	Scopes    string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type IntrospectAccessTokenModel struct {
	ClientId string
	Scopes   string
	Revoked  bool
}

type RevokeTokenModel struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
}

type RefreshTokenModel struct {
	TokenHash string
	UserId    string
	ClientId  string
	Scopes    string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}

type IntrospectRefreshTokenModel struct {
	ClientId  string
	Scopes    string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
}
