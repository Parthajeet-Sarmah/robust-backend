package controllers

import (
	"local/bomboclat-oauth-server/services"
	custom_types "local/bomboclat-oauth-server/types"
	"local/bomboclat-oauth-server/utils"
	"net/http"
)

type ClientController struct{}

func (controller ClientController) Register(w http.ResponseWriter, r *http.Request) {

	var m custom_types.ClientDatabaseModelInput

	if err := utils.DecodeJSONBody(w, r, &m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := services.ClientService.Register(&m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
