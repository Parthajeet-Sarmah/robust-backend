package main

import (
	"log"
	"net/http"

	database "local/bomboclat-oauth-server/database"
	"local/bomboclat-oauth-server/routers"
	"local/bomboclat-oauth-server/services"
	utils "local/bomboclat-oauth-server/utils"
)

func main() {

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
	authorizationRouter := routers.AuthorizationHandler().RegisterRoutes()
	userRouter := routers.UserHandler().RegisterRoutes()

	router := http.NewServeMux()

	router.Handle("/authorize/", http.StripPrefix("/authorize", authorizationRouter))
	router.Handle("/users/", http.StripPrefix("/users", userRouter))
	//router.Handle("/clients/", http.StripPrefix("/clients", clientRouter))
	//router.Handle("/introspect/", http.StripPrefix("/introspect", introspectRouter))

	log.Println("Starting server on port 9000")
	http.ListenAndServe(":9000", router)
}
