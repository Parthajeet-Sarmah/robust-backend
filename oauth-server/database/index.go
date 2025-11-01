package database

import (
	custom_types "local/bomboclat-oauth-server/types"
	"local/bomboclat-oauth-server/utils"
	"log"
)

func CreateDatabaseTables(dbPool *custom_types.Postgres) {

	tables := []struct {
		name string
		fn   func(*custom_types.Postgres) error
	}{
		{"clients", utils.CreateClientsTable},
		{"consents", utils.CreateConsentsTable},
		{"auth_codes", utils.CreateAuthCodesTable},
		{"access_tokens", utils.CreateAccessTokensTable},
		{"refresh_tokens", utils.CreateRefreshTokensTable},
	}

	for _, t := range tables {
		if err := t.fn(dbPool); err != nil {
			log.Fatalf("Could not create %s table: %v", t.name, err)
		}
	}

}
