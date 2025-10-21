package routers

import (
	controller "local/bomboclat-oauth-server/controllers"
	"net/http"
)

type UserRouterHandler struct{}

func UserHandler() *UserRouterHandler {
	return &UserRouterHandler{}
}

func (h *UserRouterHandler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	controller := controller.UserController{}

	//r.HandleFunc("POST /register", controller.Register)
	r.HandleFunc("/login", controller.Login)
	//r.HandleFunc("POST /authorize/consent", controller.AuthorizeConsent)

	return r
}
