package controllers

import (
	"local/bomboclat-oauth-server/services"
	"net/http"
)

type ClientController struct{}

//	type ClientDatabaseModelInput struct {
//		ClientSecretHash string
//		RedirectUri      string
//		AppName          string
//		GrantTypes       []string
//		LogoUrl          string
//		JwksUri          string
//		IsConfidential   bool
func (controller ClientController) Register(w http.ResponseWriter, r *http.Request) {

	services.ClientService.Register()
}
