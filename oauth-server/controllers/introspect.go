package controllers

import (
	"encoding/json"
	"net/http"

	"local/bomboclat-oauth-server/services"
	custom_types "local/bomboclat-oauth-server/types"
)

type IntrospectController struct{}

func (controller *IntrospectController) Introspect(w http.ResponseWriter, r *http.Request) {

	m := &custom_types.InstrospectModelInput{
		Token:         r.FormValue("token"),
		TokenTypeHint: r.FormValue("token_type_hint"),
	}

	metadata, err := services.IntrospectService.Introspect(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
