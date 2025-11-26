package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	custom_errors "local/bomboclat-oidc-service/errors"
	"local/bomboclat-oidc-service/services"
	custom_types "local/bomboclat-oidc-service/types"
)

type UserController struct{}

func (userController UserController) GetUserById(w http.ResponseWriter, r *http.Request) {

	user_id := r.PathValue("id")
	fields := r.URL.Query().Get("fields")

	data, err := services.UserService.GetUserById(user_id, fields)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

	randomBytes := make([]byte, 128)

	if _, err := rand.Read(randomBytes); err != nil {
		log.Print("Error while reading random bytes for generating code!")
		panic(err)
	}

	opuasValue := hex.EncodeToString(randomBytes)

	//user agent state cookie
	opuas := &http.Cookie{
		Name:     "opuas",
		Value:    opuasValue,
		Domain:   "localhost",
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)
	http.SetCookie(w, opuas)

	next := r.FormValue("next")
	if next == "" {
		next = "/"
	}

	// NOTE: Get base URL for the OAuth service
	oauthBaseUrl := os.Getenv("OAUTH_BASE_URL")
	http.Redirect(w, r, oauthBaseUrl+r.URL.Query().Get("next"), http.StatusFound)
}

func (controller UserController) Logout(w http.ResponseWriter, r *http.Request) {

	userCookie, err := r.Cookie("session_id")

	if err != nil {
		appError := custom_errors.UserNotLoggedInError(err)
		http.Error(w, appError.Error(), appError.HttpStatus)
		return
	}

	err = services.UserService.Logout(userCookie)

	if err != nil {
		var appError *custom_errors.AppError
		if errors.As(err, &appError) {
			http.Error(w, appError.Error(), appError.HttpStatus)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
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

	if !strings.HasPrefix(authHeader, "Bearer ") {
		err := custom_errors.Internal("An unexpected error occured!", errors.New("Invalid internal service token!"))
		http.Error(w, err.Error(), err.HttpStatus)
		return
	}

	authToken := strings.TrimPrefix(authHeader, "Bearer ")

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
