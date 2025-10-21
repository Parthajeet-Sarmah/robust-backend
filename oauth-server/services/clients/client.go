package clients

import (
	"context"
	"local/bomboclat-oauth-server/database"
	"local/bomboclat-oauth-server/models"
)

//type ClientDatabaseModelInput struct {
//	ClientSecretHash string
//	RedirectUri      string
//	AppName          string
//	GrantTypes       []string
//	LogoUrl          string
//	JwksUri          string
//	IsConfidential   bool

func (cs *ClientService) Register(m *models.ClientDatabaseModelInput) error {

	if err := database.InsertClient(cs.DBConn, context.Background(), m); err != nil {
		return err
	}

	return nil
}
