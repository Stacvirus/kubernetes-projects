package routes

import (
	"net/http"

	"github.com/Stacvirus/hash-generator-app/internal/handler"
)

func RegisterRoutes(mux *http.ServeMux, handler *handler.HandleRequest) {
	mux.HandleFunc("/", handler.RouteHandler)
	mux.HandleFunc("/pings", handler.CounterHandler)
}
