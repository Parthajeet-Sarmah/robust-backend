package custom_types

import "github.com/golang-jwt/jwt/v5"

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
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

type TokenModelInput struct {
	GrantType           string `json:"grant_type"`
	Code                string `json:"code"`
	RedirectUri         string `json:"redirect_uri"`
	ClientId            string `json:"client_id"`
	ClientSecretHash    string `json:"client_secret_hash,omitempty"`
	CodeVerifier        string `json:"code_verifier"`
	ClientAssertionType string `json:"client_assertion_type,omitempty"`
	ClientAssertion     string `json:"client_assertion,omitempty"`
	CodeChallengeMethod string `json:"code_challenge_method,omitempty"`
	RefreshToken        string `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	IdToken      string `json:"idToken,omitempty"`
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
