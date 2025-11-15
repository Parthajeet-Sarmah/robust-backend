package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"
	"os"

	database "local/bomboclat-oauth-server/database"
	"local/bomboclat-oauth-server/middlewares"
	"local/bomboclat-oauth-server/routers"
	"local/bomboclat-oauth-server/services"
	custom_types "local/bomboclat-oauth-server/types"
	utils "local/bomboclat-oauth-server/utils"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	//Init database and inject to services
	dbPool, err := utils.CreateDBConnPool()
	if err != nil {
		log.Fatal(err)
	}

	middlewares.InjectDBToServices(dbPool)
	services.InjectDBToServices(dbPool)
	database.CreateDatabaseTables(dbPool)

	//Init redis client and inject to services
	redisClient, err := utils.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	middlewares.InjectRedisClientToServices(redisClient)
	services.InjectRedisClientToServices(redisClient)

	//Sub routes
	authorizationRouter := routers.AuthorizationHandler().RegisterRoutes()
	clientRouter := routers.ClientHandler().RegisterRoutes()
	introspectRouter := routers.IntrospectHandler().RegisterRoutes()

	router := http.NewServeMux()

	router.Handle("/authorize/", http.StripPrefix("/authorize", authorizationRouter))
	router.Handle("/clients/", http.StripPrefix("/clients", clientRouter))
	router.Handle("/introspect/", http.StripPrefix("/introspect", introspectRouter))

	router.HandleFunc("GET /.well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		publicKey := os.Getenv("JWT_RSA_PUBLIC_KEY")

		// NOTE: Decode PEM to the original base64 encoding
		block, _ := pem.Decode([]byte(publicKey))

		if block == nil || block.Type != "PUBLIC KEY" {
			return
		}

		// NOTE: Convert the key to into crypto.PublicKey
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return
		}

		// NOTE: Cast key to rsa.PublicKey
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return
		}

		// NOTE: Get the modulus (N) and encode it to base64
		n := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())

		// NOTE: Get the exponent (E), convert it to big-endian order, and encode it to base64
		eBytes := make([]byte, 0)
		for e := rsaPub.E; e > 0; e >>= 8 {
			eBytes = append([]byte{byte(e & 0xff)}, eBytes...)
		}
		e := base64.URLEncoding.EncodeToString(eBytes)

		jwk := custom_types.JWK{
			Kty: "RSA",
			Kid: os.Getenv("JWK_KEY_ID"),
			N:   n,
			E:   e,
			Alg: "RSA256",
			Use: "sig",
		}

		jwks := map[string][]custom_types.JWK{
			"keys": []custom_types.JWK{jwk},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(jwks); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

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
			IdTokenSigningAlgValuesSupported:  []string{"RS256"},
			TokenEndpointAuthMethodsSupported: []string{"client_secret_basic"},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting server on port 9000")
	http.ListenAndServe(":9000", middlewares.CORS(router))
}
