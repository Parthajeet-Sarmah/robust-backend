package controllers

import (
	"context"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uri, err := url.Parse(m.RedirectUri)
	if err != nil || !uri.IsAbs() {
		cerr := custom_errors.InvalidRedirectURI(nil)
		http.Error(w, cerr.Error(), cerr.HttpStatus)
		return
	}
	if uri.Fragment != "" || uri.User != nil || uri.Scheme == "javascript" {
		cerr := custom_errors.InvalidRedirectURI(nil)
		http.Error(w, cerr.Error(), cerr.HttpStatus)
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
		if uri.Scheme == "http" && uri.Host != "localhost" && uri.Host != "127.0.0.1" {
			cerr := custom_errors.RedirectURIProtocolMismatch(nil)
			http.Error(w, cerr.Error(), cerr.HttpStatus)
			return
		}
	case "":
		cerr := custom_errors.Internal("An unexpected error occured!", errors.New("no environment specified"))
		http.Error(w, cerr.Error(), cerr.HttpStatus)
		return
	}

	if m.TokenEndpointAuthMethod == "private_key_jwt" {

		m.ClientSecret = ""

		if m.Jwks == "" && m.JwksUri == "" {
			cerr := custom_errors.InvalidClientMetadata(nil)
			http.Error(w, cerr.Error(), cerr.HttpStatus)
			return
		}

		if m.Jwks != "" {
			if err := utils.ValidateJwksJSON(m.Jwks); err != nil {
				cerr := custom_errors.InvalidClientMetadata(err)
				http.Error(w, cerr.Error(), cerr.HttpStatus)
				return
			}
		}

		if m.JwksUri != "" {
			juri, err := url.Parse(m.JwksUri)
			if err != nil || !juri.IsAbs() {
				cerr := custom_errors.InvalidClientMetadata(err)
				http.Error(w, cerr.Error(), cerr.HttpStatus)
				return
			}

			if env == "production" && juri.Scheme != "https" {
				cerr := custom_errors.InvalidClientMetadata(nil)
				http.Error(w, cerr.Error(), cerr.HttpStatus)
				return
			}

			if err := utils.FetchAndValidateJwks(context.Background(), m.JwksUri); err != nil {
				cerr := custom_errors.InvalidClientMetadata(err)
				http.Error(w, cerr.Error(), cerr.HttpStatus)
				return
			}
		}
	}

	if err := services.ClientService.Register(&m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
