package models

import "time"

type ClientDatabaseModelInput struct {
	ClientSecretHash string
	RedirectUri      string
	AppName          string
	GrantTypes       []string
	LogoUrl          string
	JwksUri          string
	IsConfidential   bool
}

type ClientDatabaseModel struct {
	ClientId         *string
	ClientSecretHash *string
	RedirectUri      *string
	AppName          *string
	GrantTypes       *[]string
	LogoUrl          *string
	JwksUri          *string
	IsConfidential   *bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
