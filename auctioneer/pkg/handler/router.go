package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	auctionPath  = "/auction"
	registerPath = "/register"
	listPath     = "/list"
)

func New(handler *auctionHandler) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc(auctionPath, handler.Auction).Methods(http.MethodGet)
	router.HandleFunc(registerPath, handler.Register).Methods(http.MethodPost)
	router.HandleFunc(listPath, handler.List).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 page not found", http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	})
	return router
}
