package controllers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"local/bomboclat-oidc-service/services"
	custom_types "local/bomboclat-oidc-service/types"
)

type UserController struct{}

func (userController UserController) Login(w http.ResponseWriter, r *http.Request) {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles(wd + "/templates/login.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := custom_types.UserLoginDetails{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	cookie, err := services.UserService.Login(details)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)

	next := r.FormValue("next")
	if next == "" {
		next = "/"
	}

	// NOTE: Get base URL for the OAuth service
	oauthBaseUrl := os.Getenv("OAUTH_BASE_URL")
	http.Redirect(w, r, oauthBaseUrl+r.URL.Query().Get("next"), http.StatusFound)
}

func (controller UserController) Register(w http.ResponseWriter, r *http.Request) {

	details := custom_types.UserRegistrationDetails{
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		Username:   r.FormValue("username"),
		ProfilePic: r.FormValue("profilePic"),
	}

	err := services.UserService.Register(details)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (controller UserController) UserInfo(w http.ResponseWriter, r *http.Request) {

	// TODO: Check for authorization (create a global middleware later)
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, errors.New("No auth token").Error(), http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")[1]

	if authToken == "" {
		http.Error(w, errors.New("No auth token").Error(), http.StatusUnauthorized)
		return
	}

	userInfo, err := services.UserService.UserInfo(authToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
