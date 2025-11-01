package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"local/bomboclat-oauth-server/services"
	custom_types "local/bomboclat-oauth-server/types"
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

	details := custom_types.UserDetails{
		Email:        r.FormValue("email"),
		PasswordHash: r.FormValue("password"),
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

	http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusFound)
}
