package models

import "time"

type ClientDatabaseModel struct {
	ClientId                string    `json:"client_id"`
	ClientSecretHash        *string   `json:"client_secret_hash"`         // nullable
	RedirectUri             string    `json:"redirect_uri"`               // NOT NULL
	AppName                 string    `json:"app_name"`                   // NOT NULL
	LogoUrl                 *string   `json:"logo_url"`                   // nullable
	GrantTypes              []string  `json:"grant_types"`                // NOT NULL (jsonb array)
	TokenEndpointAuthMethod string    `json:"token_endpoint_auth_method"` // NOT NULL
	Jwks                    *string   `json:"jwks"`                       // nullable
	JwksUri                 *string   `json:"jwks_uri"`                   // nullable
	IsConfidential          bool      `json:"is_confidential"`            // NOT NULL
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
