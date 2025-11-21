package custom_types

// TODO: Add multiple redirect uri support

type ClientDatabaseModelInput struct {
	ClientSecret            string   `json:"client_secret,omitempty"`
	RedirectUri             string   `json:"redirect_uri"`
	AppName                 string   `json:"app_name"`
	GrantTypes              []string `json:"grant_types"`
	LogoUrl                 string   `json:"logo_url"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"` // "client_secret_basic", "private_key_jwt"
	Jwks                    string   `json:"jwks,omitempty"`                       // optional inline JWKS JSON
	JwksUri                 string   `json:"jwks_uri,omitempty"`                   // optional
	IsConfidential          bool     `json:"is_confidential"`
}
