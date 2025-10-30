package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"local/bomboclat-oauth-server/models"
	"local/bomboclat-oauth-server/services"
	utils "local/bomboclat-oauth-server/utils"
)

type AuthorizationController struct{}

func (controller AuthorizationController) AuthorizeUserAndGenerateCode(w http.ResponseWriter, r *http.Request) {

	response_type := r.URL.Query().Get("response_type")
	client_id := r.URL.Query().Get("client_id")
	redirect_uri := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("random_state")

	//Extra security by PKCE
	code_challenge := r.URL.Query().Get("code_challenge")
	code_challenge_method := r.URL.Query().Get("code_challenge_method")

	if code_challenge_method == "" && code_challenge == "" {
		code_challenge_method = "plain"
	}

	authRequestModelInput := models.AuthorizationRequestModelInput{
		ClientId:            client_id,
		RedirectUri:         redirect_uri,
		ResponseType:        response_type,
		Scope:               scope,
		State:               state,
		CodeChallenge:       code_challenge,
		CodeChallengeMethod: code_challenge_method,
	}

	userCookie, _ := r.Cookie("session_id")

	callback_url, err := services.AuthorizationService.AuthorizeUserAndGenerateCode(authRequestModelInput, userCookie)

	if _, ok := err.(*utils.UserNotLoggedInError); ok {
		// NOTE: Redirect to /login route
		log.Println("User is not logged in")
		loginUrl := "/users/login?next=" + "/authorize" + url.QueryEscape(r.URL.String())

		http.Redirect(w, r, loginUrl, http.StatusFound)
		return
	}

	if _, ok := err.(*utils.UserScopeDeniedError); ok {

		log.Println("User has not given consent!")
		authConsentUrl := fmt.Sprintf("/authorize/consent?client_id=%s&redirect_uri=%s&scope=%s&next=%s",
			url.QueryEscape(client_id),
			url.QueryEscape(redirect_uri),
			url.QueryEscape(scope),
			"/authorize"+url.QueryEscape(r.URL.String()))

		http.Redirect(w, r, authConsentUrl, http.StatusFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, *callback_url, http.StatusFound)
}

func (controller *AuthorizationController) AuthorizeConsent(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		client_id := r.URL.Query().Get("client_id")
		scope := r.URL.Query().Get("scope")
		redirect_uri := r.URL.Query().Get("redirect_uri")
		next := r.URL.Query().Get("next")

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(template.ParseFiles(wd + "/templates/consent.html"))

		tmpl.Execute(w, struct {
			ClientId    string
			Scope       string
			RedirectUri string
			Next        string
		}{client_id, scope, redirect_uri, next})
		return
	}

	if r.Method == http.MethodPost {

		client_id := r.FormValue("client_id")
		scope := r.FormValue("scope")
		decision := r.FormValue("decision")
		redirect_uri := r.FormValue("redirect_uri")

		authConsentModelInput := models.AuthorizationConsentModelInput{
			ClientId:    client_id,
			Scope:       scope,
			Decision:    decision,
			RedirectUri: redirect_uri,
		}

		userCookie, err := r.Cookie("session_id")

		if err != nil {
			log.Print("No cookie found!")
			return
		}

		err = services.AuthorizationService.AuthorizeConsent(authConsentModelInput, userCookie)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}

		http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (controller *AuthorizationController) GenerateToken(w http.ResponseWriter, r *http.Request) {

	// TODO: Enforce proper Content-Type headers
	m := &models.TokenModelInput{
		GrantType:           r.FormValue("grant_type"),
		Code:                r.FormValue("code"),
		RedirectUri:         r.FormValue("redirect_uri"),
		ClientId:            r.FormValue("client_id"),
		ClientSecretHash:    r.FormValue("client_secret_hash"),
		CodeVerifier:        r.FormValue("code_verifier"),
		CodeChallengeMethod: r.FormValue("code_challenge_method"),
		RefreshToken:        r.FormValue("refresh_token"),
	}

	token, err := services.AuthorizationService.GenerateToken(m)

	if err != nil {

		if _, ok := err.(*utils.RefreshTokenNotFoundError); ok {
			_err := err.(*utils.RefreshTokenNotFoundError)
			http.Error(w, _err.Msg, _err.Status)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (controller *AuthorizationController) RevokeToken(w http.ResponseWriter, r *http.Request) {

	m := &models.RevokeTokenModel{
		Token:         r.FormValue("token"),
		TokenTypeHint: r.FormValue("token_type_hint"),
	}

	if err := services.AuthorizationService.RevokeToken(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
