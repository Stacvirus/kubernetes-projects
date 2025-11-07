package routes

import (
	"net/http"
	"ping-pong/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, handler *handlers.HandleRequest) {
	mux.HandleFunc("/", handler.RouteHandler)
	mux.HandleFunc("/pings", handler.CountHandler)
}
