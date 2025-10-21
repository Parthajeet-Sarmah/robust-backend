package routers

import (
	controller "local/bomboclat-oauth-server/controllers"
	"net/http"
)

type AuthorizationRouterHandler struct{}

func AuthorizationHandler() *AuthorizationRouterHandler {
	return &AuthorizationRouterHandler{}
}

func (h *AuthorizationRouterHandler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	controller := controller.AuthorizationController{}

	//r.HandleFunc("POST /register", controller.Register)
	r.HandleFunc("GET /", controller.AuthorizeUserAndGenerateCode)
	r.HandleFunc("GET /consent", controller.AuthorizeConsent)
	r.HandleFunc("POST /consent", controller.AuthorizeConsent)
	//r.HandleFunc("POST /authorize/consent", controller.AuthorizeConsent)

	return r
}
