package routers

import (
	controller "local/bomboclat-oauth-server/controllers"
	"net/http"
)

type InstrospectRouterHandler struct{}

func IntrospectHandler() *InstrospectRouterHandler {
	return &InstrospectRouterHandler{}
}

func (h *InstrospectRouterHandler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	controller := controller.IntrospectController{}

	r.HandleFunc("POST /", controller.Introspect)

	return r
}
