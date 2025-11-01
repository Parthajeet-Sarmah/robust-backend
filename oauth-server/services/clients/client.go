package clients

import (
	"context"
	"local/bomboclat-oauth-server/database"
	custom_types "local/bomboclat-oauth-server/types"
)

func (cs *ClientService) Register(m *custom_types.ClientDatabaseModelInput) error {

	if err := database.InsertClient(cs.DBConn, context.Background(), m); err != nil {
		return err
	}

	return nil
}
