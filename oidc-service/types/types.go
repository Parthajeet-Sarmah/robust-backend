package custom_types

import (
	"github.com/jackc/pgx/v5/pgxpool"

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

type OpenIdConfiguration struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserInfoEndpoint                  string   `json:"userinfo_endpoint"`
	JwksUri                           string   `json:"jwks_uri"`
	RegistrationEndpoint              string   `json:"registration_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	IntrospectionEndpoint             string   `json:"introspection_endpoint"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
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
	Email           string
	Name            string
	PasswordHash    string
	IsEmailVerified bool
	ProfilePic      string
}

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Postgres struct {
	DB *pgxpool.Pool
}
