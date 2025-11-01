package database

import (
	custom_types "local/bomboclat-oidc-service/types"
	"local/bomboclat-oidc-service/utils"
	"log"
)

func CreateDatabaseTables(dbPool *custom_types.Postgres) {

	tables := []struct {
		name string
		fn   func(*custom_types.Postgres) error
	}{
		{"users", utils.CreateUsersTable},
	}

	for _, t := range tables {
		if err := t.fn(dbPool); err != nil {
			log.Fatalf("Could not create %s table: %v", t.name, err)
		}
	}

}
