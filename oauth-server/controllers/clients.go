package controllers

import (
	"errors"
	custom_errors "local/bomboclat-oauth-server/errors"
	"local/bomboclat-oauth-server/services"
	custom_types "local/bomboclat-oauth-server/types"
	"local/bomboclat-oauth-server/utils"
	"net/http"
	"net/url"
	"os"
)

type ClientController struct{}

func (controller ClientController) Register(w http.ResponseWriter, r *http.Request) {

	var m custom_types.ClientDatabaseModelInput

	if err := utils.DecodeJSONBody(w, r, &m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	uri, err := url.Parse(m.RedirectUri)

	if err != nil {
		cerror := custom_errors.Internal("An unexpected error occured!", err)
		http.Error(w, cerror.Error(), cerror.HttpStatus)
		return
	}

	env := os.Getenv("ENVIRONMENT")

	switch env {
	case "production":
		if uri.Scheme != "https" {
			cerr := custom_errors.RedirectURIProtocolMismatch(nil)
			http.Error(w, cerr.Error(), cerr.HttpStatus)
			return
		}
	case "development":
		if uri.Host != "localhost" && uri.Host != "127.0.0.1" && uri.Scheme == "http" {
			cerr := custom_errors.RedirectURIProtocolMismatch(nil)
			http.Error(w, cerr.Error(), cerr.HttpStatus)
			return
		}
	case "":
		cerr := custom_errors.Internal("An unexpected error occured!", errors.New("No environment specified!"))
		http.Error(w, cerr.Error(), cerr.HttpStatus)
		return
	}

	if err := services.ClientService.Register(&m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
