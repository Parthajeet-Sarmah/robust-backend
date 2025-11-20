package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"local/bomboclat-oauth-server/database"
	custom_errors "local/bomboclat-oauth-server/errors"
	"local/bomboclat-oauth-server/middlewares"
	"local/bomboclat-oauth-server/models"
	"local/bomboclat-oauth-server/services"
	custom_types "local/bomboclat-oauth-server/types"
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
	codeChallenge := r.URL.Query().Get("code_challenge")
	codeChallengeMethod := r.URL.Query().Get("code_challenge_method")

	if codeChallengeMethod == "" || codeChallenge == "" {
		http.Error(w, errors.New("No code challenge method or code challenge provided").Error(), http.StatusInternalServerError)
		return
	}

	authRequestModelInput := custom_types.AuthorizationRequestModelInput{
		ClientId:            client_id,
		RedirectUri:         redirect_uri,
		ResponseType:        response_type,
		Scope:               scope,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}

	userCookie, err := r.Cookie("session_id")

	callback_url, err := services.AuthorizationService.AuthorizeUserAndGenerateCode(authRequestModelInput, userCookie)

	var appError *custom_errors.AppError

	if errors.As(err, &appError) && appError.Code == "USER_NOT_LOGGED_IN" {
		// NOTE: Redirect to /login route
		log.Println("User is not logged in")

		loginBaseUrl := os.Getenv("OIDC_BASE_URL")
		loginUrl := "/users/login?next=" + "/authorize" + url.QueryEscape(r.URL.String())

		http.Redirect(w, r, loginBaseUrl+loginUrl, http.StatusFound)
		return

	} else if errors.As(err, &appError) && appError.Code == "USER_SCOPE_DENIED" {

		log.Println("User has not given consent!")
		authConsentUrl := fmt.Sprintf("/authorize/consent?client_id=%s&redirect_uri=%s&scope=%s&next=%s",
			url.QueryEscape(client_id),
			url.QueryEscape(redirect_uri),
			url.QueryEscape(scope),
			"/authorize"+url.QueryEscape(r.URL.String()))

		http.Redirect(w, r, authConsentUrl, http.StatusFound)
		return

	} else if err != nil {
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

		// Get user session to check for existing consent
		userCookie, err := r.Cookie("session_id")
		if err != nil {
			// No session - redirect to login
			loginBaseUrl := os.Getenv("OIDC_BASE_URL")
			loginUrl := "/users/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginBaseUrl+loginUrl, http.StatusFound)
			return
		}

		// Get user ID from session
		userKey := fmt.Sprintf("user_session:%s", userCookie.Value)
		res, err := utils.GetValueFromHash(services.AuthorizationService.RedisClient, userKey)
		if err != nil || res == nil {
			// Invalid session - redirect to login
			loginBaseUrl := os.Getenv("OIDC_BASE_URL")
			loginUrl := "/users/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginBaseUrl+loginUrl, http.StatusFound)
			return
		}

		userId := res["user_id"]
		if userId == "" {
			http.Error(w, "User ID not found in session", http.StatusUnauthorized)
			return
		}

		// Check if user has already consented to this client
		consent, err := database.FindConsent(services.AuthorizationService.DBConn, userId, client_id)
		if err != nil {
			log.Printf("Error checking consent: %v", err)
			// Continue to show consent screen on error
		}

		if consent != nil {
			// User has previously consented
			// Check if requested scopes are a subset of previously granted scopes
			if scopesAreSubset(scope, consent.Scopes) {
				// Auto-approve: set consent in session and redirect to next
				res["scope"] = "allow"
				utils.SetValueToHash(services.AuthorizationService.RedisClient, userKey, res)

				if next == "" {
					next = "/"
				}
				http.Redirect(w, r, next, http.StatusFound)
				return
			}
			// If requesting additional scopes, show consent screen
		}

		// Show consent screen (first time or requesting new scopes)
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
		next := r.FormValue("next")

		authConsentModelInput := custom_types.AuthorizationConsentModelInput{
			ClientId:    client_id,
			Scope:       scope,
			Decision:    decision,
			RedirectUri: redirect_uri,
		}

		userCookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "No session cookie found", http.StatusUnauthorized)
			return
		}

		// Get user ID from session
		userKey := fmt.Sprintf("user_session:%s", userCookie.Value)
		res, err := utils.GetValueFromHash(services.AuthorizationService.RedisClient, userKey)
		if err != nil || res == nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		userId := res["user_id"]
		if userId == "" {
			http.Error(w, "User ID not found in session", http.StatusUnauthorized)
			return
		}

		err = services.AuthorizationService.AuthorizeConsent(authConsentModelInput, userCookie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if decision == "allow" {
			err = database.UpsertConsent(services.AuthorizationService.DBConn, &models.ConsentModel{
				UserId:   userId,
				ClientId: client_id,
				Scopes:   scope,
			})
			if err != nil {
				log.Printf("Error saving consent: %v", err)
			}
		}

		if next == "" {
			next = "/"
		}

		http.Redirect(w, r, next, http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (controller *AuthorizationController) GenerateToken(w http.ResponseWriter, r *http.Request) {

	// TODO: Enforce proper Content-Type headers
	m := &custom_types.TokenModelInput{
		GrantType:           r.FormValue("grant_type"),
		Code:                r.FormValue("code"),
		RedirectUri:         r.FormValue("redirect_uri"),
		ClientId:            r.FormValue("client_id"),
		ClientSecretHash:    r.FormValue("client_secret_hash"),
		CodeVerifier:        r.FormValue("code_verifier"),
		CodeChallengeMethod: r.FormValue("code_challenge_method"),
		ClientAssertion:     r.FormValue("client_assertion"),
		ClientAssertionType: r.FormValue("client_assertion_type"),
		RefreshToken:        r.FormValue("refresh_token"),
	}

	// NOTE: Middleware to check if client is authorized (client_secret_basic)
	_, err := middlewares.MiddlewareService.AuthorizeClient(r)

	authMethod := "client_secret_basic"

	if err != nil {

		// NOTE: Check for other authentication types (client_secret_post, private_key_jwt)
		clientSecretPost := m != &custom_types.TokenModelInput{} && m.ClientSecretHash != ""
		privateKeyJwt := m.ClientAssertionType != "" && m.ClientAssertion != ""
		none := m.CodeVerifier != ""

		hasOnlyOneMethod := (clientSecretPost && !privateKeyJwt && !none) ||
			(!clientSecretPost && privateKeyJwt && !none) ||
			(!clientSecretPost && !privateKeyJwt && none)

		if !hasOnlyOneMethod {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if clientSecretPost {
			authMethod = "client_secret_post"
		} else if privateKeyJwt {
			authMethod = "private_key_jwt"
		} else if none {
			authMethod = "none"
		}
	}

	token, err := services.AuthorizationService.GenerateToken(m, authMethod)

	if err != nil {
		var appErr *custom_errors.AppError

		if errors.As(err, &appErr) {
			http.Error(w, appErr.Error(), appErr.HttpStatus)
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

// scopesAreSubset checks if requestedScopes is a subset of grantedScopes
func scopesAreSubset(requestedScopes, grantedScopes string) bool {
	// Split space-separated scopes
	requested := strings.Split(strings.TrimSpace(requestedScopes), " ")
	granted := strings.Split(strings.TrimSpace(grantedScopes), " ")

	// Create a map of granted scopes for quick lookup
	grantedMap := make(map[string]bool)
	for _, scope := range granted {
		if scope != "" {
			grantedMap[scope] = true
		}
	}

	// Check if all requested scopes are in granted scopes
	for _, scope := range requested {
		if scope != "" && !grantedMap[scope] {
			return false // Requesting a new scope
		}
	}

	return true
}
