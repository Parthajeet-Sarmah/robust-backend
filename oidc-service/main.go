package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	database "local/bomboclat-oidc-service/database"
	"local/bomboclat-oidc-service/routers"
	"local/bomboclat-oidc-service/services"
	custom_types "local/bomboclat-oidc-service/types"
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

	//Configuration endpoints
	router.HandleFunc("GET /.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {

		config := custom_types.OpenIdConfiguration{
			Issuer:                            os.Getenv("OIDC_BASE_URL"),
			UserInfoEndpoint:                  os.Getenv("OIDC_BASE_URL") + "userinfo",
			JwksUri:                           os.Getenv("OIDC_BASE_URL") + ".well-known/jwks.json",
			AuthorizationEndpoint:             os.Getenv("OAUTH_BASE_URL") + "authorize",
			TokenEndpoint:                     os.Getenv("OAUTH_BASE_URL") + "authorize/token",
			RegistrationEndpoint:              os.Getenv("OAUTH_BASE_URL") + "clients/register",
			RevocationEndpoint:                os.Getenv("OAUTH_BASE_URL") + "authorize/revoke",
			IntrospectionEndpoint:             os.Getenv("OAUTH_BASE_URL") + "introspect",
			ScopesSupported:                   []string{"openid", "profile", "email"},
			ResponseTypesSupported:            []string{"code"},
			GrantTypesSupported:               []string{"authorization_code", "refresh_token"},
			SubjectTypesSupported:             []string{"public"},
			IdTokenSigningAlgValuesSupported:  []string{"S256"},
			TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"},
		}

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.HandleFunc("GET /.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {

	})

	log.Println("Starting server on port 9030")
	http.ListenAndServe(":9030", router)
}
