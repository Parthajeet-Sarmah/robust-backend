package controllers

import (
	"errors"
	"html/template"
	custom_errors "local/bomboclat-oidc-service/errors"
	"local/bomboclat-oidc-service/services"
	custom_types "local/bomboclat-oidc-service/types"
	"log"
	"net/http"
	"os"
	"time"
)

type SessionController struct{}

func (controller SessionController) CheckSession(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(template.ParseFiles(wd + "/templates/login-status.html"))

	tmpl.Execute(w, nil)
}

func (controller SessionController) EndSession(w http.ResponseWriter, r *http.Request) {

	userCookie, err := r.Cookie("session_id")

	idTokenHint := r.URL.Query().Get("id_token_hint")
	postLogoutRedirectUri := r.URL.Query().Get("post_logout_redirect_uri")

	m := custom_types.EndSessionInput{
		UserCookie:            userCookie,
		IdTokenHint:           idTokenHint,
		PostLogoutRedirectUri: postLogoutRedirectUri,
	}

	if err != nil {
		appError := custom_errors.UserNotLoggedInError(err)
		http.Error(w, appError.Error(), appError.HttpStatus)
		return
	}

	err = services.SessionService.EndSession(m)

	if err != nil {
		var appError *custom_errors.AppError
		if errors.As(err, &appError) {
			http.Error(w, appError.Error(), appError.HttpStatus)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	//session id cookie removal
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Secure:   true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	}

	//user agent state cookie removal
	opuas := &http.Cookie{
		Name:     "opuas",
		Value:    "",
		Domain:   "localhost",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	}

	http.SetCookie(w, opuas)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
