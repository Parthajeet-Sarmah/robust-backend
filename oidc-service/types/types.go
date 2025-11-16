package custom_types

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-jwt/jwt/v5"
)

type IntrospectionResponse struct {
	ClientID  string   `json:"client_id"`
	Scope     string   `json:"scope"`
	Active    bool     `json:"active"`
	Sub       string   `json:"sub"`
	Exp       int64    `json:"exp"`
	Iat       int64    `json:"iat"`
	Iss       string   `json:"iss"`
	TokenType string   `json:"token_type"`
	Aud       []string `json:"aud"`
}

type CustomClaims struct {
	Issuer    string `json:"iss"`
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	NotBefore int64  `json:"nbf"`
	IssuedAt  int64  `json:"iat"`
	jwt.RegisteredClaims
}

type UserLoginDetails struct {
	Email    string
	Password string
}

type UserRegistrationDetails struct {
	Username        string `json:"username" db:"username"`
	Email           string `json:"email" db:"email"`
	Password        string `json:"password" db:"password"`
	IsEmailVerified bool   `json:"email_verified" db:"email_verified"`
	ProfilePic      string `json:"profile_pic" db:"profile_pic"`
}

type UserProfile struct {
	UserUUID        string `json:"id"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	PasswordHash    string `json:"password_hash,omitempty"`
	IsEmailVerified bool   `json:"email_verified,omitempty"`
	ProfilePic      string `json:"profile_pic,omitempty"`
}

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Postgres struct {
	DB *pgxpool.Pool
}
