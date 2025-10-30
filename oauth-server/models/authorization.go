package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	NotBefore int64  `json:"nbf"`
	IssuedAt  int64  `json:"iat"`
	jwt.RegisteredClaims
}

type AuthorizationRequestModelInput struct {
	ResponseType        string
	ClientId            string
	RedirectUri         string
	Scope               string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type AuthorizationConsentModelInput struct {
	ClientId    string
	Scope       string
	Decision    string
	RedirectUri string
}

type AuthCodeModelInput struct {
	Id                  string
	Code                string
	UserId              string
	ClientId            string
	RedirectUri         string
	Scopes              string
	CodeChallenge       string
	CodeChallengeMethod string
	Used                bool
}

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

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
}

type TokenModelInput struct {
	GrantType           string `json:"grant_type"`
	Code                string `json:"code"`
	RedirectUri         string `json:"redirect_uri"`
	ClientId            string `json:"client_id"`
	ClientSecretHash    string `json:"client_secret_hash"`
	CodeVerifier        string `json:"code_verifier"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	RefreshToken        string `json:"refresh_token"`
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

type InstrospectModelInput struct {
	Token         string
	TokenTypeHint string
}
