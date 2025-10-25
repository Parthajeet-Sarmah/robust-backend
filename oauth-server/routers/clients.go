package routers

import (
	controller "local/bomboclat-oauth-server/controllers"
	"net/http"
)

type ClientRouterHandler struct{}

func ClientHandler() *ClientRouterHandler {
	return &ClientRouterHandler{}
}

func (h *ClientRouterHandler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	controller := controller.ClientController{}

	r.HandleFunc("POST /register", controller.Register)

	return r
}
