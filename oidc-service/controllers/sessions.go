package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
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

}
