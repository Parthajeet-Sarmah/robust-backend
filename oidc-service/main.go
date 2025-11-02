package main

import (
	"log"
	"net/http"

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

	router := http.NewServeMux()

	router.Handle("/users/", http.StripPrefix("/users", userRouter))

	log.Println("Starting server on port 9030")
	http.ListenAndServe(":9030", router)
}
