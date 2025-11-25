package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	database "local/bomboclat-oidc-service/database"
	"local/bomboclat-oidc-service/routers"
	"local/bomboclat-oidc-service/services"
	utils "local/bomboclat-oidc-service/utils"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	//Init database and inject to services
	dbPool, err := utils.CreateDBConnPool()
	if err != nil {
		log.Fatal(err)
	}
	services.InjectDBToServices(dbPool)
	database.CreateDatabaseTables(dbPool)

	//Init redis client and inject to services
	redisClient, err := utils.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
		return
	}
	services.InjectRedisClientToServices(redisClient)

	//Sub routes
	userRouter := routers.UserHandler().RegisterRoutes()
	sessionsRouter := routers.SessionHandler().RegisterRoutes()

	router := http.NewServeMux()

	router.HandleFunc("/login-status", func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		tmpl := template.Must(template.ParseFiles(wd + "/templates/login-status.html"))

		if r.Method == http.MethodPost {
			tmpl.Execute(w, nil)
		}
	})

	router.Handle("/sessions/", http.StripPrefix("/sessions", sessionsRouter))
	router.Handle("/users/", http.StripPrefix("/users", userRouter))

	log.Println("Starting server on port 9030")
	http.ListenAndServe(":9030", router)
}
