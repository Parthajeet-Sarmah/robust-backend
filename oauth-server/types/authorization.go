package custom_types

import "github.com/golang-jwt/jwt/v5"

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

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
}

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
