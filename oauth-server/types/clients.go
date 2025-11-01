package custom_types

type ClientDatabaseModelInput struct {
	ClientSecretHash string   `json:"client_secret_hash"`
	RedirectUri      string   `json:"redirect_uri"`
	AppName          string   `json:"app_name"`
	GrantTypes       []string `json:"grant_types"`
	LogoUrl          string   `json:"logo_url"`
	JwksUri          string   `json:"jwks_uri"`
	IsConfidential   bool     `json:"is_confidential"`
}
