package routers

import (
	controller "local/bomboclat-oidc-service/controllers"
	"net/http"
)

type SessionRouterHandler struct{}

func SessionHandler() *SessionRouterHandler {
	return &SessionRouterHandler{}
}

func (h *SessionRouterHandler) RegisterRoutes() *http.ServeMux {
	r := http.NewServeMux()

	controller := controller.SessionController{}

	r.HandleFunc("GET /check", controller.CheckSession)
	r.HandleFunc("GET /end_session", controller.EndSession)

	return r
}
